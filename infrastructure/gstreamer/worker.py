#!/usr/bin/env python3
"""
GStreamer Worker Pool - Issue #24
Consumes transcoding jobs from Kafka and processes them using GStreamer
"""

import os
import json
import subprocess
import logging
import time
from kafka import KafkaConsumer
from kafka.errors import KafkaError
import redis
import requests

# Configuration
KAFKA_BROKERS = os.getenv('KAFKA_BROKERS', 'localhost:9092').split(',')
KAFKA_TOPIC = os.getenv('KAFKA_TOPIC', 'transcoding')
REDIS_HOST = os.getenv('REDIS_HOST', 'localhost')
REDIS_PORT = int(os.getenv('REDIS_PORT', 6379))
TRANSCODING_SERVICE_URL = os.getenv('TRANSCODING_SERVICE_URL', 'http://transcoding-service:8080')
GPU_ENABLED = os.getenv('GPU_ENABLED', 'false').lower() == 'true'

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Initialize Redis for job checkpointing
redis_client = redis.Redis(host=REDIS_HOST, port=REDIS_PORT, decode_responses=True)


def build_gstreamer_pipeline(job):
    """
    Build GStreamer pipeline for transcoding job
    
    GPU acceleration (NVIDIA):
    - nvh264enc (H.264)
    - nvh265enc (H.265)
    
    CPU fallback:
    - x264enc (H.264)
    - x265enc (H.265)
    """
    input_url = job['input_url']
    output_url = job['output_url']
    codec = job.get('codec', 'h264')
    width = job.get('width', 1920)
    height = job.get('height', 1080)
    bitrate = job.get('bitrate', 5000000)
    fps = job.get('fps', 30)
    
    # Determine encoder based on codec and GPU availability
    if GPU_ENABLED and codec == 'h264':
        encoder = 'nvh264enc'
        encoder_options = f'preset=fast bitrate={bitrate}'
    elif GPU_ENABLED and codec == 'h265':
        encoder = 'nvh265enc'
        encoder_options = f'preset=fast bitrate={bitrate}'
    elif codec == 'h264':
        encoder = 'x264enc'
        encoder_options = f'speed-preset=fast bitrate={bitrate}'
    elif codec == 'h265':
        encoder = 'x265enc'
        encoder_options = f'speed-preset=fast bitrate={bitrate}'
    else:
        encoder = 'x264enc'
        encoder_options = f'speed-preset=fast bitrate={bitrate}'
    
    # Build pipeline
    # filesrc location=input.mp4 ! qtdemux ! avdec_h264 ! videoscale ! x264enc ! h264parse ! mpegtsmux ! filesink location=output.ts
    pipeline = f"""
    gst-launch-1.0 -e \
    filesrc location={input_url} ! \
    qtdemux name=demux \
    demux.video_0 ! queue ! decodebin ! \
    videoscale ! video/x-raw,width={width},height={height} ! \
    {encoder} {encoder_options} ! \
    {codec}parse ! \
    mpegtsmux ! \
    filesink location={output_url}
    """
    
    return pipeline.strip()


def process_transcode_job(job):
    """Process a transcoding job"""
    job_id = job['job_id']
    
    try:
        # Checkpoint: Start processing
        checkpoint_key = f"transcode:checkpoint:{job_id}"
        redis_client.set(checkpoint_key, json.dumps({
            'status': 'processing',
            'progress': 0,
            'started_at': time.time()
        }))
        
        # Update job status to processing
        requests.patch(
            f"{TRANSCODING_SERVICE_URL}/transcode/jobs/{job_id}/progress",
            json={'progress': 0, 'status': 'processing'}
        )
        
        # Build pipeline
        pipeline = build_gstreamer_pipeline(job)
        logger.info(f"Building pipeline for job {job_id}: {pipeline[:100]}...")
        
        # Execute GStreamer pipeline
        process = subprocess.Popen(
            pipeline,
            shell=True,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True
        )
        
        # Monitor progress (simplified - TODO: parse actual progress from GStreamer)
        start_time = time.time()
        while process.poll() is None:
            elapsed = time.time() - start_time
            # Estimate progress (TODO: get actual from GStreamer)
            estimated_duration = job.get('duration', 3600)  # Assume 1 hour if unknown
            progress = min(90, int((elapsed / estimated_duration) * 100))
            
            redis_client.set(checkpoint_key, json.dumps({
                'status': 'processing',
                'progress': progress,
                'elapsed': elapsed
            }))
            
            requests.patch(
                f"{TRANSCODING_SERVICE_URL}/transcode/jobs/{job_id}/progress",
                json={'progress': progress}
            )
            
            time.sleep(5)
        
        stdout, stderr = process.communicate()
        
        if process.returncode == 0:
            # Success
            redis_client.set(checkpoint_key, json.dumps({
                'status': 'completed',
                'progress': 100,
                'completed_at': time.time()
            }))
            
            requests.patch(
                f"{TRANSCODING_SERVICE_URL}/transcode/jobs/{job_id}/complete",
                json={'output_url': job['output_url'], 'status': 'completed'}
            )
            
            logger.info(f"Job {job_id} completed successfully")
        else:
            # Failure
            error_msg = stderr or "Transcoding failed"
            redis_client.set(checkpoint_key, json.dumps({
                'status': 'failed',
                'error': error_msg,
                'failed_at': time.time()
            }))
            
            requests.patch(
                f"{TRANSCODING_SERVICE_URL}/transcode/jobs/{job_id}/fail",
                json={'error': error_msg, 'status': 'failed'}
            )
            
            logger.error(f"Job {job_id} failed: {error_msg}")
    
    except Exception as e:
        logger.error(f"Error processing job {job_id}: {str(e)}")
        requests.patch(
            f"{TRANSCODING_SERVICE_URL}/transcode/jobs/{job_id}/fail",
            json={'error': str(e), 'status': 'failed'}
        )


def main():
    """Main worker loop"""
    logger.info("Starting GStreamer worker...")
    logger.info(f"Kafka brokers: {KAFKA_BROKERS}")
    logger.info(f"Kafka topic: {KAFKA_TOPIC}")
    logger.info(f"GPU enabled: {GPU_ENABLED}")
    
    # Create Kafka consumer
    consumer = KafkaConsumer(
        KAFKA_TOPIC,
        bootstrap_servers=KAFKA_BROKERS,
        value_deserializer=lambda m: json.loads(m.decode('utf-8')),
        group_id='gstreamer-workers',
        auto_offset_reset='earliest',
        enable_auto_commit=True,
    )
    
    logger.info("Worker started. Waiting for jobs...")
    
    try:
        for message in consumer:
            job = message.value
            logger.info(f"Received job: {job['job_id']}")
            process_transcode_job(job)
    
    except KeyboardInterrupt:
        logger.info("Worker shutting down...")
    except KafkaError as e:
        logger.error(f"Kafka error: {str(e)}")
    finally:
        consumer.close()


if __name__ == '__main__':
    main()

