FROM nvidia/cuda:11.8.0-cudnn8-runtime-ubuntu22.04

# Install dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    pkg-config \
    python3 \
    python3-pip \
    libgstreamer1.0-dev \
    libgstreamer-plugins-base1.0-dev \
    libgstreamer-plugins-bad1.0-dev \
    gstreamer1.0-plugins-base \
    gstreamer1.0-plugins-good \
    gstreamer1.0-plugins-bad \
    gstreamer1.0-plugins-ugly \
    gstreamer1.0-libav \
    gstreamer1.0-tools \
    gstreamer1.0-x \
    gstreamer1.0-alsa \
    gstreamer1.0-gl \
    gstreamer1.0-gtk3 \
    gstreamer1.0-qt5 \
    gstreamer1.0-pulseaudio \
    libsoup2.4-dev \
    libjpeg-dev \
    libpng-dev \
    zlib1g-dev \
    libgstreamer1.0-dev \
    libgstreamer-plugins-base1.0-dev \
    && rm -rf /var/lib/apt/lists/*

# Install NVIDIA GPU plugins for GStreamer
RUN apt-get update && apt-get install -y \
    gstreamer1.0-nvcodec \
    gstreamer1.0-nvcodec-dev \
    || echo "GPU plugins not available, continuing without GPU acceleration"

# Install Kafka Python client
RUN pip3 install kafka-python redis requests

# Create worker directory
WORKDIR /app

# Copy worker script
COPY worker.py .
COPY transcode_job.py .

# Set environment
ENV PYTHONUNBUFFERED=1

# Run worker
CMD ["python3", "worker.py"]

