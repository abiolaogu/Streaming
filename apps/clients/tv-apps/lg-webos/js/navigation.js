/**
 * TV Navigation Handler for webOS Magic Remote
 */

class WebOSNavigation {
    constructor() {
        this.currentFocus = null;
        this.focusableElements = [];
        this.init();
    }

    init() {
        document.addEventListener('keydown', (e) => {
            this.handleKey(e);
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
        const keyCode = e.keyCode;
        const keyName = e.key;

        // webOS Magic Remote keys
        switch(keyCode) {
            case 38: case 65362: // Up
                e.preventDefault();
                this.moveFocus('up');
                break;
            case 40: case 65364: // Down
                e.preventDefault();
                this.moveFocus('down');
                break;
            case 37: case 65361: // Left
                e.preventDefault();
                this.moveFocus('left');
                break;
            case 39: case 65363: // Right
                e.preventDefault();
                this.moveFocus('right');
                break;
            case 13: case 65376: // Enter
                e.preventDefault();
                if (this.currentFocus) {
                    this.currentFocus.click();
                }
                break;
            case 27: case 65385: // Back/Escape
                e.preventDefault();
                this.handleBack();
                break;
            case 10009: // Back (webOS)
                e.preventDefault();
                this.handleBack();
                break;
        }

        // Pointer mode (Magic Remote pointer)
        if (keyName === 'PointerMove') {
            // Handle pointer movement if needed
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
                window.location.hash = '#detail';
            }
        }
    }
}

const webosNavigation = new WebOSNavigation();

