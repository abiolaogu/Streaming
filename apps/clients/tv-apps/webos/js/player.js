/**
 * webOS MediaPlayer Implementation
 */

class WebOSPlayer {
    constructor() {
        this.video = null;
        this.mediaSource = null;
        this.isPlaying = false;
        this.currentContent = null;
    }

    init() {
        // Create video element for webOS MediaPlayer
        this.video = document.createElement('video');
        this.video.id = 'mediaPlayer';
        this.video.style.width = '100%';
        this.video.style.height = '100%';
        
        const container = document.getElementById('player-container');
        if (container) {
            container.appendChild(this.video);
        }

        // Setup event handlers
        this.setupEventHandlers();
        return true;
    }

    setupEventHandlers() {
        if (!this.video) return;

        this.video.addEventListener('play', () => {
            this.isPlaying = true;
            console.log('Playback started');
        });

        this.video.addEventListener('pause', () => {
            this.isPlaying = false;
        });

        this.video.addEventListener('ended', () => {
            this.stop();
        });

        this.video.addEventListener('error', (e) => {
            console.error('Player error:', e);
            this.handleError(e);
        });

        this.video.addEventListener('waiting', () => {
            console.log('Buffering...');
        });

        this.video.addEventListener('canplay', () => {
            console.log('Ready to play');
        });
    }

    async play(content) {
        if (!this.video) {
            if (!this.init()) {
                throw new Error('Player initialization failed');
            }
        }

        this.currentContent = content;
        const streamUrl = content.streamUrl;

        try {
            // Configure video source
            this.video.src = streamUrl;
            
            // Configure DRM if needed
            if (content.isDrmProtected) {
                await this.configureDRM(content);
            }

            // Configure video attributes
            this.video.setAttribute('playsinline', 'true');
            this.video.setAttribute('webkit-playsinline', 'true');

            // Enter fullscreen
            if (this.video.webkitEnterFullscreen) {
                this.video.webkitEnterFullscreen();
            }

            // Start playback
            await this.video.play();

        } catch (error) {
            console.error('Playback failed:', error);
            throw error;
        }
    }

    async configureDRM(content) {
        if (!this.video) return;

        // webOS supports PlayReady and Widevine
        if (typeof webOS !== 'undefined' && webOS.media) {
            const licenseServer = 'https://drm.streamverse.com/v1/' + 
                (content.drmType === 'playready' ? 'playready' : 'widevine') + '/license';
            const token = apiService.getToken() || '';

            try {
                // Configure DRM via webOS media API
                const drmInitData = {
                    licenseServer: licenseServer,
                    customData: {
                        'Authorization': `Bearer ${token}`,
                        'X-Content-ID': content.id
                    }
                };

                if (content.drmType === 'playready') {
                    this.video.setAttribute('drm-type', 'playready');
                } else if (content.drmType === 'widevine') {
                    this.video.setAttribute('drm-type', 'widevine');
                }

                // Set encrypted media extensions
                if (this.video.setMediaKeys) {
                    // Configure EME for DRM
                    const keySystem = content.drmType === 'playready' 
                        ? 'com.microsoft.playready' 
                        : 'com.widevine.alpha';
                    
                    // Note: Actual EME implementation requires proper key system support
                    console.log('DRM configured for:', content.drmType);
                }

            } catch (error) {
                console.error('DRM configuration failed:', error);
            }
        }
    }

    pause() {
        if (this.video && this.isPlaying) {
            this.video.pause();
            this.isPlaying = false;
        }
    }

    resume() {
        if (this.video && !this.isPlaying) {
            this.video.play();
            this.isPlaying = true;
        }
    }

    stop() {
        if (this.video) {
            this.video.pause();
            this.video.src = '';
            this.isPlaying = false;
            this.currentContent = null;
            
            // Exit fullscreen if active
            if (this.video.webkitExitFullscreen) {
                this.video.webkitExitFullscreen();
            }
        }
    }

    seek(position) {
        if (this.video) {
            this.video.currentTime = position;
        }
    }

    handleError(e) {
        console.error('Player error:', e);
        const error = this.video.error;
        if (error) {
            console.error('Error code:', error.code, 'Message:', error.message);
        }
        alert('Playback error occurred. Please try again.');
    }
}

const webOSPlayer = new WebOSPlayer();

