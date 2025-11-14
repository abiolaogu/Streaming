class AppConfig {
  // Environment
  static const String environment = String.fromEnvironment(
    'ENVIRONMENT',
    defaultValue: 'production',
  );

  // API Configuration
  static const String baseUrl = String.fromEnvironment(
    'API_BASE_URL',
    defaultValue: 'https://api.streamverse.io',
  );

  static const String apiVersion = 'v1';
  static String get apiUrl => '$baseUrl/api/$apiVersion';

  // Service Endpoints
  static String get authServiceUrl => '$apiUrl/auth';
  static String get userServiceUrl => '$apiUrl/users';
  static String get contentServiceUrl => '$apiUrl/content';
  static String get streamingServiceUrl => '$apiUrl/streaming';
  static String get paymentServiceUrl => '$apiUrl/payments';
  static String get searchServiceUrl => '$apiUrl/search';
  static String get recommendationServiceUrl => '$apiUrl/recommendations';

  // API Keys
  static const String geminiApiKey = String.fromEnvironment('GEMINI_API_KEY');
  static const String sentryDsn = String.fromEnvironment('SENTRY_DSN');

  // App Configuration
  static const String appName = 'StreamVerse';
  static const String appVersion = '1.0.0';
  static const int appBuildNumber = 1;

  // Feature Flags
  static const bool enableOfflineDownloads = true;
  static const bool enableChromecast = true;
  static const bool enablePictureInPicture = true;
  static const bool enableBiometricAuth = true;
  static const bool enableSocialFeatures = true;

  // Media Configuration
  static const int maxDownloadQuality = 1080;
  static const int defaultStreamingQuality = 720;
  static const bool autoPlayNextEpisode = true;
  static const int autoPlayCountdownSeconds = 5;

  // Cache Configuration
  static const Duration cacheDuration = Duration(hours: 24);
  static const int maxCacheSize = 500; // MB

  // Network Configuration
  static const Duration connectionTimeout = Duration(seconds: 30);
  static const Duration receiveTimeout = Duration(seconds: 30);
  static const int maxRetries = 3;

  // Pagination
  static const int itemsPerPage = 20;

  // DRM Configuration
  static const String widevineLicenseUrl =
      'https://drm.streamverse.io/widevine';
  static const String fairPlayLicenseUrl =
      'https://drm.streamverse.io/fairplay';

  // CDN URLs
  static const String cdnBaseUrl = 'https://cdn.streamverse.io';
  static String getThumbnailUrl(String contentId) =>
      '$cdnBaseUrl/thumbnails/$contentId';
  static String getPosterUrl(String contentId) =>
      '$cdnBaseUrl/posters/$contentId';

  // Analytics
  static const bool enableAnalytics = true;
  static const bool enableCrashReporting = true;
}
