# Ad Decision & Placement Logic (ADPL)

This service is the core logic unit for the ad-compositing platform. It receives requests, gathers necessary data, and makes the final decision on which ad to show and where to place it.

## Core Technologies

*   **Language:** Go

## Key Responsibilities

-   Act as the main API endpoint for the `ad-compositing-service`.
-   Fetch scene analysis data from the VSAE.
-   Fetch user data from the `user-service`.
-   Request ad creatives from the `ad-inventory-service`.
-   Apply business rules and logic to select the optimal ad and placement.
-   Send final instructions to the Real-time Compositing Engine (RCE).
