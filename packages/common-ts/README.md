# Common TypeScript Package

Shared TypeScript utilities and types for StreamVerse frontend applications.

## Features

- ✅ API client with axios
- ✅ Type definitions
- ✅ Storage utilities
- ✅ Token management

## Usage

```typescript
import { ApiClient, Content, TokenStorage } from '@streamverse/common-ts';

const apiClient = new ApiClient({ baseURL: 'https://api.streamverse.com' });

// Set token
TokenStorage.setAccessToken('token-here');
apiClient.setToken(TokenStorage.getAccessToken());

// Make API calls
const content = await apiClient.get<Content>('/api/v1/content/123');
```

