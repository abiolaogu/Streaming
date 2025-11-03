/**
 * LG webOS MediaPlayer Implementation
 */

class WebOSPlayer {
    constructor() {
        this.player = null;
        this.isPlaying = false;
        this.currentContent = null;
    }

    init() {
        try {
            // Create MediaPlayer
            this.player = new MediaPlayer();
            this.setupEventHandlers();
            return true;
        } catch (error) {
            console.error('Failed to initialize player:', error);
            return false;
        }
    }

    setupEventHandlers() {
        if (!this.player) return;

        this.player.addEventListener('loadstatechange', (event) => {
            if (event.loadState === 'loaded') {
                console.log('Media loaded');
            }
        });

        this.player.addEventListener('playstatechange', (event) => {
            if (event.playState === 'playing') {
                this.isPlaying = true;
                this.updatePlayerUI();
            } else if (event.playState === 'paused') {
                this.isPlaying = false;
            }
        });

        this.player.addEventListener('bufferingstart', () => {
            this.showBuffering();
        });

        this.player.addEventListener('bufferingcomplete', () => {
            this.hideBuffering();
        });

        this.player.addEventListener('error', (error) => {
            console.error('Player error:', error);
            this.handleError(error);
        });
    }

    async play(content) {
        if (!this.player) {
            if (!this.init()) {
                throw new Error('Player initialization failed');
            }
        }

        this.currentContent = content;
        const streamUrl = content.streamUrl;

        try {
            // Create video element
            const video = document.createElement('video');
            video.id = 'video-player';
            video.style.width = '100%';
            video.style.height = '100%';

            const container = document.getElementById('player-container');
            container.innerHTML = '';
            container.appendChild(video);

            // Configure DRM if needed
            if (content.isDrmProtected) {
                await this.configureDRM(video, content);
            }

            video.src = streamUrl;
            video.play();

            this.player = video;

        } catch (error) {
            console.error('Playback failed:', error);
            throw error;
        }
    }

    async configureDRM(video, content) {
        const licenseServer = 'https://drm.streamverse.com/v1/widevine/license';
        const token = apiService.getToken() || '';

        // Configure EME (Encrypted Media Extensions)
        try {
            const keySystem = content.drmType === 'playready' 
                ? 'com.microsoft.playready' 
                : 'com.widevine.alpha';

            video.addEventListener('encrypted', async (event) => {
                try {
                    const keySystemAccess = await navigator.requestMediaKeySystemAccess(keySystem, [{
                        initDataTypes: ['cenc'],
                        audioCapabilities: [{
                            contentType: 'audio/mp4;codecs="mp4a.40.2"',
                            robustness: 'SW_SECURE_CRYPTO'
                        }],
                        videoCapabilities: [{
                            contentType: 'video/mp4;codecs="avc1.42E01E"',
                            robustness: 'SW_SECURE_CRYPTO'
                        }]
                    }]);

                    const mediaKeys = await keySystemAccess.createMediaKeys();
                    await video.setMediaKeys(mediaKeys);
                    const session = mediaKeys.createSession();
                    
                    session.addEventListener('message', async (event) => {
                        const response = await fetch(licenseServer, {
                            method: 'POST',
                            headers: {
                                'Content-Type': 'application/octet-binary',
                                'Authorization': `Bearer ${token}`,
                                'X-Content-ID': content.id
                            },
                            body: event.message
                        });
                        
                        const license = await response.arrayBuffer();
                        session.update(license);
                    });

                    await session.generateRequest('cenc', event.initData);
                } catch (error) {
                    console.error('DRM configuration failed:', error);
                }
            });
        } catch (error) {
            console.error('DRM setup error:', error);
        }
    }

    pause() {
        if (this.player && this.isPlaying) {
            this.player.pause();
            this.isPlaying = false;
        }
    }

    resume() {
        if (this.player && !this.isPlaying) {
            this.player.play();
            this.isPlaying = true;
        }
    }

    stop() {
        if (this.player) {
            this.player.pause();
            this.player.src = '';
            this.isPlaying = false;
            this.currentContent = null;
            
            const container = document.getElementById('player-container');
            container.innerHTML = '';
        }
    }

    updatePlayerUI() {
        const playerScreen = document.getElementById('player-screen');
        if (playerScreen) {
            playerScreen.classList.add('active');
        }
    }

    showBuffering() {
        console.log('Showing buffering indicator');
    }

    hideBuffering() {
        console.log('Hiding buffering indicator');
    }

    handleError(error) {
        console.error('Player error:', error);
        alert('Playback error occurred. Please try again.');
    }
}

const webosPlayer = new WebOSPlayer();

