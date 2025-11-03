/**
 * webOS Magic Remote Navigation Handler
 */

class WebOSNavigation {
    constructor() {
        this.currentFocus = null;
        this.focusableElements = [];
        this.init();
    }

    init() {
        // Listen for key events (Magic Remote)
        document.addEventListener('keydown', (e) => {
            this.handleKey(e);
        });

        // Listen for pointer events (Magic Remote pointer)
        document.addEventListener('pointermove', (e) => {
            this.handlePointer(e);
        });

        this.updateFocusableElements();
        this.focusFirst();
    }

    updateFocusableElements() {
        this.focusableElements = Array.from(
            document.querySelectorAll('button, input, select, .content-card, [tabindex="0"]')
        ).filter(el => {
            const style = window.getComputedStyle(el);
            return style.display !== 'none' && style.visibility !== 'hidden';
        });
    }

    handleKey(e) {
        const keyCode = e.keyCode || e.which;

        switch(keyCode) {
            case 38: // Up
                e.preventDefault();
                this.moveFocus('up');
                break;
            case 40: // Down
                e.preventDefault();
                this.moveFocus('down');
                break;
            case 37: // Left
                e.preventDefault();
                this.moveFocus('left');
                break;
            case 39: // Right
                e.preventDefault();
                this.moveFocus('right');
                break;
            case 13: // Enter
                e.preventDefault();
                if (this.currentFocus) {
                    this.currentFocus.click();
                }
                break;
            case 461: // Back (webOS)
            case 27: // Escape
                e.preventDefault();
                this.handleBack();
                break;
        }
    }

    handlePointer(e) {
        // Magic Remote pointer support
        const element = document.elementFromPoint(e.clientX, e.clientY);
        if (element && this.focusableElements.includes(element)) {
            this.setFocus(element);
        }
    }

    moveFocus(direction) {
        if (!this.currentFocus) {
            this.focusFirst();
            return;
        }

        const currentIndex = this.focusableElements.indexOf(this.currentFocus);
        let nextIndex = currentIndex;

        if (direction === 'left' || direction === 'right') {
            nextIndex = direction === 'left' ? currentIndex - 1 : currentIndex + 1;
        } else {
            const rect = this.currentFocus.getBoundingClientRect();
            const centerX = rect.left + rect.width / 2;
            
            if (direction === 'down') {
                nextIndex = this.findNextInColumn(centerX, rect.bottom);
            } else {
                nextIndex = this.findPreviousInColumn(centerX, rect.top);
            }
        }

        if (nextIndex >= 0 && nextIndex < this.focusableElements.length) {
            this.setFocus(this.focusableElements[nextIndex]);
        }
    }

    findNextInColumn(x, y) {
        let bestIndex = -1;
        let bestDistance = Infinity;

        this.focusableElements.forEach((el, index) => {
            const rect = el.getBoundingClientRect();
            const elX = rect.left + rect.width / 2;
            const distance = Math.abs(x - elX);

            if (rect.top > y && distance < bestDistance) {
                bestDistance = distance;
                bestIndex = index;
            }
        });

        return bestIndex >= 0 ? bestIndex : this.focusableElements.length - 1;
    }

    findPreviousInColumn(x, y) {
        let bestIndex = -1;
        let bestDistance = Infinity;

        this.focusableElements.forEach((el, index) => {
            const rect = el.getBoundingClientRect();
            const elX = rect.left + rect.width / 2;
            const distance = Math.abs(x - elX);

            if (rect.top < y && distance < bestDistance) {
                bestDistance = distance;
                bestIndex = index;
            }
        });

        return bestIndex >= 0 ? bestIndex : 0;
    }

    setFocus(element) {
        if (this.currentFocus) {
            this.currentFocus.classList.remove('focused');
        }
        this.currentFocus = element;
        element.classList.add('focused');
        element.focus();
    }

    focusFirst() {
        if (this.focusableElements.length > 0) {
            this.setFocus(this.focusableElements[0]);
        }
    }

    handleBack() {
        const activeScreen = document.querySelector('.screen.active');
        if (activeScreen) {
            const screenId = activeScreen.id;
            
            if (screenId === 'detail-screen' || screenId === 'search-screen' || screenId === 'settings-screen') {
                window.location.hash = '#home';
            } else if (screenId === 'player-screen') {
                webOSPlayer.stop();
                window.location.hash = '#detail';
            }
        }
    }
}

const webOSNavigation = new WebOSNavigation();

