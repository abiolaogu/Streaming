package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	runpodAPIBaseURL = "https://api.runpod.io"
	minQueueDepth    = 5  // Trigger at 5 pending jobs
	maxInstances     = 20 // Max RunPod instances
	idleThreshold    = 60 // Seconds without work before scale-down
)

type Config struct {
	RunPodAPIKey     string
	InstanceType     string
	NetworkVolumeId  string
	TemplateId       string
	KubernetesMaster string
	KubernetesCA     string
}

type RunPodPodsResponse struct {
	Pods []RunPodPod `json:"pods"`
}

type RunPodPod struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	MachineId    string `json:"machineId"`
	MachineType  string `json:"machineType"`
	Status       string `json:"desiredStatus"`
	PublicIp     string `json:"publicIp"`
	RuntimeInUse bool   `json:"runtimeInUse"`
}

type KafkaQueueDepth struct {
	Topic  string
	Depth  int64
	Rate   float64
}

type KedaQueueMetric struct {
	Results []struct {
		Metric map[string]string `json:"metric"`
		Value  []interface{}     `json:"value"`
	} `json:"result"`
}

func main() {
	log.Println("RunPod Autoscaler starting...")

	config := loadConfig()

	clientset, err := getKubernetesClient()
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	lastScaleTime := time.Now()
	var mu sync.Mutex

	for {
		select {
		case <-ticker.C:
			mu.Lock()
			
			queueDepth, err := getKafkaQueueDepth(ctx)
			if err != nil {
				log.Printf("Error getting queue depth: %v", err)
				continue
			}

			runningPods, err := listRunPodInstances(config)
			if err != nil {
				log.Printf("Error listing RunPod instances: %v", err)
				continue
			}

			log.Printf("Queue depth: %d, Running instances: %d", queueDepth, len(runningPods))

			// Scale-up decision
			if queueDepth >= minQueueDepth && len(runningPods) < maxInstances {
				needed := min(queueDepth/10+1, maxInstances-len(runningPods))
				log.Printf("Scaling up %d RunPod instances", needed)
				
				for i := 0; i < needed; i++ {
					pod, err := createRunPodInstance(config)
					if err != nil {
						log.Printf("Error creating RunPod instance: %v", err)
						continue
					}
					
					log.Printf("Created RunPod instance: %s at %s", pod.Id, pod.PublicIp)
					
					// Attach to Kubernetes cluster via SSH
					go attachToKubernetes(pod, config)
				}
				
				lastScaleTime = time.Now()
			}

			// Scale-down decision
			if queueDepth == 0 && len(runningPods) > 0 && time.Since(lastScaleTime) > time.Duration(idleThreshold)*time.Second {
				log.Printf("Scaling down RunPod instances (idle)")
				
				for _, pod := range runningPods {
					if err := deleteRunPodInstance(config, pod.Id); err != nil {
						log.Printf("Error deleting pod %s: %v", pod.Id, err)
					} else {
						log.Printf("Deleted RunPod instance: %s", pod.Id)
					}
				}
				
				lastScaleTime = time.Now()
			}

			mu.Unlock()
		case <-ctx.Done():
			log.Println("Shutting down...")
			return
		}
	}
}

func loadConfig() Config {
	return Config{
		RunPodAPIKey:     getEnv("RUNPOD_API_KEY", ""),
		InstanceType:     getEnv("RUNPOD_INSTANCE_TYPE", "NVIDIA RTX 6000 Ada Generation"),
		NetworkVolumeId:  getEnv("RUNPOD_NETWORK_VOLUME_ID", ""),
		TemplateId:       getEnv("RUNPOD_TEMPLATE_ID", ""),
		KubernetesMaster: getEnv("KUBERNETES_MASTER", ""),
		KubernetesCA:     getEnv("KUBERNETES_CA", ""),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getKubernetesClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	
	clientset, err := kubernetes.NewForConfig(config)
	return clientset, err
}

func getKafkaQueueDepth(ctx context.Context) (int64, error) {
	// Query KEDA metrics endpoint for queue depth
	url := "http://keda-metrics-apiserver.keda:8080/api/v1/query?query=keda_scaler_metrics_value{scalerName=\"kafka-scaler\",namespace=\"data\"}"
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, err
	}
	
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var metric KedaQueueMetric
	if err := json.NewDecoder(resp.Body).Decode(&metric); err != nil {
		return 0, err
	}

	if len(metric.Results) > 0 && len(metric.Results[0].Value) >= 2 {
		if depthStr, ok := metric.Results[0].Value[1].(string); ok {
			var depth int64
			fmt.Sscanf(depthStr, "%d", &depth)
			return depth, nil
		}
	}

	return 0, nil
}

func listRunPodInstances(config Config) ([]RunPodPod, error) {
	url := fmt.Sprintf("%s/pods", runpodAPIBaseURL)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Authorization", "Bearer "+config.RunPodAPIKey)
	
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var podsResp RunPodPodsResponse
	if err := json.NewDecoder(resp.Body).Decode(&podsResp); err != nil {
		return nil, err
	}

	// Filter active pods
	var active []RunPodPod
	for _, pod := range podsResp.Pods {
		if pod.Status == "RUNNING" || pod.Status == "READY" {
			active = append(active, pod)
		}
	}

	return active, nil
}

func createRunPodInstance(config Config) (*RunPodPod, error) {
	payload := map[string]interface{}{
		"machineType":    config.InstanceType,
		"networkVolumeId": config.NetworkVolumeId,
		"templateId":     config.TemplateId,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/pods", runpodAPIBaseURL)
	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+config.RunPodAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pod RunPodPod
	if err := json.NewDecoder(resp.Body).Decode(&pod); err != nil {
		return nil, err
	}

	return &pod, nil
}

func deleteRunPodInstance(config Config, podId string) error {
	url := fmt.Sprintf("%s/pods/%s", runpodAPIBaseURL, podId)
	
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+config.RunPodAPIKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func attachToKubernetes(pod *RunPodPod, config Config) {
	log.Printf("Attaching RunPod instance %s to Kubernetes", pod.Id)
	
	// SSH into the RunPod instance and install kubelet
	sshCmd := exec.Command("ssh", 
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		fmt.Sprintf("root@%s", pod.PublicIp),
		"bash", "-c", fmt.Sprintf(`
			# Install Kubernetes node
			curl -fsSL https://get.k8s.io | KUBERNETES_VERSION=1.28 bash -
			
			# Configure kubelet to join cluster
			echo 'KUBELET_EXTRA_ARGS="--node-ip=%s"' > /etc/default/kubelet
			
			# Join cluster using CA and master
			kubeadm join %s --discovery-token-unsafe-skip-ca-verification
		`, pod.PublicIp, config.KubernetesMaster))
	
	sshCmd.Env = append(os.Environ(), fmt.Sprintf("KUBERNETES_CA=%s", config.KubernetesCA))
	
	output, err := sshCmd.CombinedOutput()
	if err != nil {
		log.Printf("Error attaching to Kubernetes: %v\nOutput: %s", err, output)
		return
	}
	
	log.Printf("Successfully attached RunPod instance %s to Kubernetes", pod.Id)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

