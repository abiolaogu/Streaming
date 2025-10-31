/*
 * STB Cache Daemon for Satellite NIP Reception
 * Listens on DVB-S2X multicast and serves HTTP to local network
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <pthread.h>
#include <sys/mman.h>
#include <fcntl.h>
#include <errno.h>

#define MULTICAST_GROUP "239.255.1.1"
#define MULTICAST_PORT 5001
#define HTTP_PORT 8080
#define CACHE_DIR "/var/cache/satellite"
#define MAX_CACHE_SIZE (10ULL * 1024 * 1024 * 1024) // 10GB

typedef struct {
    char id[64];
    char url[512];
    size_t size;
    time_t expiry;
    void* data;
    pthread_rwlock_t lock;
} cache_entry_t;

typedef struct {
    cache_entry_t** entries;
    size_t count;
    size_t capacity;
    pthread_rwlock_t lock;
} cache_db_t;

cache_db_t cache_db;

static void* multicast_receiver(void* arg) {
    int sock = socket(AF_INET, SOCK_DGRAM, 0);
    if (sock < 0) {
        perror("socket");
        return NULL;
    }

    struct sockaddr_in addr;
    memset(&addr, 0, sizeof(addr));
    addr.sin_family = AF_INET;
    addr.sin_addr.s_addr = INADDR_ANY;
    addr.sin_port = htons(MULTICAST_PORT);

    if (bind(sock, (struct sockaddr*)&addr, sizeof(addr)) < 0) {
        perror("bind");
        close(sock);
        return NULL;
    }

    struct ip_mreq mreq;
    mreq.imr_multiaddr.s_addr = inet_addr(MULTICAST_GROUP);
    mreq.imr_interface.s_addr = INADDR_ANY;

    if (setsockopt(sock, IPPROTO_IP, IP_ADD_MEMBERSHIP, &mreq, sizeof(mreq)) < 0) {
        perror("setsockopt IP_ADD_MEMBERSHIP");
        close(sock);
        return NULL;
    }

    printf("Listening on %s:%d for NIP carousel\n", MULTICAST_GROUP, MULTICAST_PORT);

    char buffer[188]; // DVB TS packet size
    while (1) {
        ssize_t len = recvfrom(sock, buffer, sizeof(buffer), 0, NULL, NULL);
        if (len < 0) {
            perror("recvfrom");
            continue;
        }

        // Process TS packet and extract carousel data
        // This is simplified - real implementation would parse DSM-CC Object Carousel
        
        // TODO: Parse TS, extract DSM-CC files, update cache
        printf("Received TS packet: %zd bytes\n", len);
    }

    close(sock);
    return NULL;
}

static void* http_server(void* arg) {
    int sock = socket(AF_INET, SOCK_STREAM, 0);
    if (sock < 0) {
        perror("socket");
        return NULL;
    }

    int opt = 1;
    setsockopt(sock, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt));

    struct sockaddr_in addr;
    memset(&addr, 0, sizeof(addr));
    addr.sin_family = AF_INET;
    addr.sin_addr.s_addr = INADDR_ANY;
    addr.sin_port = htons(HTTP_PORT);

    if (bind(sock, (struct sockaddr*)&addr, sizeof(addr)) < 0) {
        perror("bind");
        close(sock);
        return NULL;
    }

    if (listen(sock, 10) < 0) {
        perror("listen");
        close(sock);
        return NULL;
    }

    printf("HTTP server listening on port %d\n", HTTP_PORT);

    while (1) {
        int client = accept(sock, NULL, NULL);
        if (client < 0) {
            perror("accept");
            continue;
        }

        char request[1024];
        ssize_t len = recv(client, request, sizeof(request) - 1, 0);
        if (len < 0) {
            close(client);
            continue;
        }

        request[len] = '\0';

        // Simple HTTP request parser
        char method[16], path[256];
        sscanf(request, "%s %s", method, path);

        printf("Request: %s %s\n", method, path);

        // Lookup in cache or serve from terrestrial CDN
        const char* response = "HTTP/1.1 200 OK\r\n"
                              "Content-Type: video/mp4\r\n"
                              "Cache-Control: public, max-age=3600\r\n"
                              "Connection: close\r\n\r\n";

        send(client, response, strlen(response), 0);

        // TODO: Send actual content from cache or proxy to CDN

        close(client);
    }

    close(sock);
    return NULL;
}

static void cleanup_expired_entries() {
    pthread_rwlock_wrlock(&cache_db.lock);

    time_t now = time(NULL);
    for (size_t i = 0; i < cache_db.count; i++) {
        if (cache_db.entries[i] && cache_db.entries[i]->expiry < now) {
            printf("Expiring cache entry: %s\n", cache_db.entries[i]->id);
            munmap(cache_db.entries[i]->data, cache_db.entries[i]->size);
            pthread_rwlock_destroy(&cache_db.entries[i]->lock);
            free(cache_db.entries[i]);
            cache_db.entries[i] = NULL;
        }
    }

    pthread_rwlock_unlock(&cache_db.lock);
}

int main(int argc, char* argv[]) {
    printf("STB Cache Daemon starting...\n");

    // Initialize cache database
    cache_db.capacity = 1000;
    cache_db.count = 0;
    cache_db.entries = calloc(cache_db.capacity, sizeof(cache_entry_t*));
    pthread_rwlock_init(&cache_db.lock, NULL);

    // Create cache directory
    char cmd[256];
    snprintf(cmd, sizeof(cmd), "mkdir -p %s", CACHE_DIR);
    system(cmd);

    // Start multicast receiver thread
    pthread_t receiver_thread;
    if (pthread_create(&receiver_thread, NULL, multicast_receiver, NULL) != 0) {
        perror("pthread_create receiver");
        return 1;
    }

    // Start HTTP server thread
    pthread_t http_thread;
    if (pthread_create(&http_thread, NULL, http_server, NULL) != 0) {
        perror("pthread_create http");
        return 1;
    }

    printf("STB Cache Daemon running\n");

    // Periodic cleanup
    while (1) {
        sleep(300); // Cleanup every 5 minutes
        cleanup_expired_entries();
    }

    pthread_join(receiver_thread, NULL);
    pthread_join(http_thread, NULL);
    pthread_rwlock_destroy(&cache_db.lock);

    return 0;
}

