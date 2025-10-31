CC=gcc
CFLAGS=-Wall -Wextra -std=c99 -pthread -O2
TARGET=stb-cache-daemon
SOURCES=stb-cache-daemon.c

.PHONY: all clean install

all: $(TARGET)

$(TARGET): $(SOURCES)
	$(CC) $(CFLAGS) -o $(TARGET) $(SOURCES)

clean:
	rm -f $(TARGET)

install: $(TARGET)
	sudo mkdir -p /var/cache/satellite
	sudo cp $(TARGET) /usr/local/bin/
	sudo cp stb-cache-daemon.service /etc/systemd/system/
	sudo systemctl daemon-reload
	sudo systemctl enable stb-cache-daemon
	sudo systemctl start stb-cache-daemon

