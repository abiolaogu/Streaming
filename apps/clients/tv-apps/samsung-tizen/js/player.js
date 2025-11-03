/**
 * Tizen AVPlay Player Implementation
 */

class TizenPlayer {
    constructor() {
        this.player = null;
        this.isPlaying = false;
        this.currentContent = null;
    }

    init() {
        try {
            // Create AVPlay player
            this.player = new webapis.avplay();
            this.setupEventHandlers();
            return true;
        } catch (error) {
            console.error('Failed to initialize player:', error);
            return false;
        }
    }

    setupEventHandlers() {
        if (!this.player) return;

        this.player.setListener({
            onstreamstarted: () => {
                console.log('Stream started');
                this.isPlaying = true;
                this.updatePlayerUI();
            },
            onstreamcompleted: () => {
                console.log('Stream completed');
                this.stop();
            },
            onerror: (eventType) => {
                console.error('Player error:', eventType);
                this.handleError(eventType);
            },
            onbufferingstart: () => {
                console.log('Buffering started');
                this.showBuffering();
            },
            onbufferingprogress: (percent) => {
                console.log('Buffering progress:', percent);
            },
            onbufferingcomplete: () => {
                console.log('Buffering complete');
                this.hideBuffering();
            },
            oncurrentplaytime: (currentTime) => {
                this.updatePlayTime(currentTime);
            }
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
            // Set stream URL
            this.player.open(streamUrl);

            // Configure DRM if needed
            if (content.isDrmProtected && content.drmType === 'playready') {
                await this.configureDRM(content);
            }

            // Set display mode
            this.player.setDisplayRect(0, 0, 1920, 1080);

            // Prepare and play
            this.player.prepareAsync(() => {
                this.player.play();
            });

        } catch (error) {
            console.error('Playback failed:', error);
            throw error;
        }
    }

    async configureDRM(content) {
        if (!this.player) return;

        const licenseServer = 'https://drm.streamverse.com/v1/playready/license';
        const token = apiService.getToken() || '';

        // Configure PlayReady DRM
        const drmInitData = {
            licenseServer: licenseServer,
            customData: JSON.stringify({
                'Authorization': `Bearer ${token}`,
                'X-Content-ID': content.id
            })
        };

        this.player.setDrm(webapis.avplay.PLAYREADY, drmInitData);
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
            this.player.stop();
            this.player.close();
            this.isPlaying = false;
            this.currentContent = null;
        }
    }

    seek(position) {
        if (this.player) {
            this.player.seekTo(position);
        }
    }

    updatePlayerUI() {
        // Update player UI elements
        const playerScreen = document.getElementById('player-screen');
        if (playerScreen) {
            playerScreen.classList.add('active');
        }
    }

    showBuffering() {
        // Show buffering indicator
        console.log('Showing buffering indicator');
    }

    hideBuffering() {
        // Hide buffering indicator
        console.log('Hiding buffering indicator');
    }

    updatePlayTime(currentTime) {
        // Update play time display
        console.log('Current time:', currentTime);
    }

    handleError(eventType) {
        console.error('Player error type:', eventType);
        // Show error message to user
        alert('Playback error occurred. Please try again.');
    }
}

const tizenPlayer = new TizenPlayer();

