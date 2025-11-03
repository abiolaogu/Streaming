# Payment Service

Microservice for subscription management and billing.

## Features

- ✅ Subscription plans (Basic, Premium)
- ✅ Subscription management (subscribe, cancel, pause)
- ✅ TVOD (Transactional VOD) - Rent/Buy
- ✅ PPV (Pay-Per-View) support
- ✅ Payment processing
- ✅ Invoice generation

## API Endpoints

- `POST /api/v1/payments/subscribe` - Subscribe to plan
- `GET /api/v1/payments/subscription` - Get subscription
- `POST /api/v1/payments/subscription/cancel` - Cancel subscription
- `POST /api/v1/payments/purchase` - Create purchase (rent/buy)

## Running

```bash
go run main.go
```

