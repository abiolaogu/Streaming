/**
 * Local storage utilities
 */

export class Storage {
  private static prefix = 'streamverse_';

  static setItem(key: string, value: string): void {
    if (typeof window !== 'undefined' && window.localStorage) {
      window.localStorage.setItem(this.prefix + key, value);
    }
  }

  static getItem(key: string): string | null {
    if (typeof window !== 'undefined' && window.localStorage) {
      return window.localStorage.getItem(this.prefix + key);
    }
    return null;
  }

  static removeItem(key: string): void {
    if (typeof window !== 'undefined' && window.localStorage) {
      window.localStorage.removeItem(this.prefix + key);
    }
  }

  static clear(): void {
    if (typeof window !== 'undefined' && window.localStorage) {
      const keys = Object.keys(window.localStorage);
      keys.forEach(key => {
        if (key.startsWith(this.prefix)) {
          window.localStorage.removeItem(key);
        }
      });
    }
  }
}

// Token management
export const TokenStorage = {
  setAccessToken: (token: string) => Storage.setItem('access_token', token),
  getAccessToken: (): string | null => Storage.getItem('access_token'),
  setRefreshToken: (token: string) => Storage.setItem('refresh_token', token),
  getRefreshToken: (): string | null => Storage.getItem('refresh_token'),
  clear: () => {
    Storage.removeItem('access_token');
    Storage.removeItem('refresh_token');
  },
};

