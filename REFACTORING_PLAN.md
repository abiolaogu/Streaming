# Refactoring Plan: Streaming Service

This document outlines the refactoring plan for the `streaming-service`. The goal is to align the service with the new microservices architecture, integrate it with other services, and implement the missing features.

## 1. Service Integration

*   **Integrate with Content Service**: Replace the mock `ContentRepository` with a gRPC client that communicates with the `content-service`.
*   **Integrate with Payment Service**: Replace the mock `SubscriptionRepository` with a gRPC client that communicates with the `payment-service`.

## 2. Implement QoE Analytics

*   **Integrate with Kafka**: Implement a Kafka producer to send QoE events to the `qoe-events` topic.
*   **Define QoE Event Schema**: Define a clear and consistent schema for QoE events.

## 3. Implement Core Features

*   **Device Detection**: Implement a mechanism to detect the user's device type and select the optimal ABR profile.
*   **Geo-Restrictions**: Implement a geo-restriction check to ensure that users are only able to access content from allowed regions.

## 4. Code Cleanup

*   **Remove TODOs**: Address all `// TODO:` comments in the codebase.
*   **Improve Error Handling**: Implement more robust and consistent error handling.
*   **Add Unit and Integration Tests**: Write comprehensive unit and integration tests for all new features and refactored code.

## 5. Timeline

This refactoring effort will be completed in a series of small, incremental steps. Each step will be a separate pull request, which will be reviewed and tested before being merged into the main branch.
