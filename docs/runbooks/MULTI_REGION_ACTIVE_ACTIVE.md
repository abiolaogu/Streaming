# Multi-Region Active/Active Runbook

## Purpose

Operational runbook for global active/active routing, regional data residency controls, and Stripe reconciliation workers.

## Topology Baseline

- Regions: `us-east-1`, `us-west-2`, `eu-west-1` (extend as needed)
- Global edge routing: latency-based DNS or global load balancer
- Service deployment: active in every region
- Data strategy:
  - User/auth/payment data pinned to residency-allowed region(s)
  - Read-mostly catalog/content metadata replicated globally
  - Event processing local-first with asynchronous cross-region replication

## Required Environment Controls

- `PLATFORM_REGION`: current runtime region (for each deployment)
- `ACTIVE_ACTIVE_REGIONS`: comma-separated routable regions
- `DATA_RESIDENCY_MODE`: `global` or `strict`
- `DATA_RESIDENCY_ALLOWED_REGIONS`: regions allowed to process residency-bound events
- `STRIPE_RECON_WORKER_ENABLED`: enable failed-event replay worker
- `STRIPE_RECON_WORKER_INTERVAL`: replay cadence (for example `60s`)
- `STRIPE_RECON_WORKER_BATCH_SIZE`: replay batch limit

## Regional Residency Guardrails

1. Set `DATA_RESIDENCY_MODE=strict` for production.
2. Set `DATA_RESIDENCY_ALLOWED_REGIONS` to jurisdiction-approved regions.
3. Ensure each deployment sets `PLATFORM_REGION` correctly.
4. In strict mode, reconciliation workers in disallowed regions must stay disabled automatically.

## Active/Active Traffic Policy

1. Health checks must pass in all active regions:
   - `/health` for auth/content/payment/policy services.
2. Global router should prefer lowest-latency healthy region.
3. On regional degradation, route away from unhealthy region.
4. Keep at least one fallback region warm for every tenant geography.

## Stripe Reconciliation Worker Operations

The payment service runs a reconciliation worker that replays failed webhook events from durable idempotency storage.

### Normal checks

1. Verify worker enabled:
   - service log contains `Stripe reconciliation worker` startup messages.
2. Verify replay progress:
   - failed event backlog trends downward.
3. Confirm idempotency:
   - replayed events move from `failed` to `processed` without duplicate mutations.

### Backlog response

1. Temporarily increase `STRIPE_RECON_WORKER_BATCH_SIZE`.
2. Decrease `STRIPE_RECON_WORKER_INTERVAL` (for example `30s`).
3. Scale payment-service replicas in residency-allowed regions only.
4. Monitor DB write pressure and error rates.

## Regional Failover Procedure

1. Detect regional outage via health/error alerts.
2. Remove failing region from global routing set.
3. Confirm traffic drains to remaining active regions.
4. Validate auth, entitlement checks, playback, and payment webhook ingestion.
5. Re-add recovered region after smoke checks pass.

## Validation Checklist

- `./scripts/ci/validate-monorepo.sh all`
- Coverage gates and contract gates pass
- Policy-service and payment-service health endpoints healthy in each region
- No unresolved failed webhook backlog in primary residency regions
