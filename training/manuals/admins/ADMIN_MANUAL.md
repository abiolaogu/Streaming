# StreamVerse Administrator Manual

**Version 2.0** | **For Platform Administrators**

---

## Table of Contents

1. [Admin Dashboard](#admin-dashboard)
2. [User Management](#user-management)
3. [Content Moderation](#content-moderation)
4. [System Configuration](#system-configuration)
5. [Analytics & Reporting](#analytics--reporting)
6. [Security & Compliance](#security--compliance)
7. [Troubleshooting](#troubleshooting)

---

## Admin Dashboard

### Accessing Admin Panel

URL: `https://admin.streamverse.io`

**Login Requirements:**
- Admin role assigned
- Multi-factor authentication (MFA) enabled
- VPN connection (for production)

### Dashboard Overview

**Key Metrics:**
- Total Users
- Active Subscriptions
- Content Library Size
- System Health Status
- Revenue (24h, 7d, 30d)
- Concurrent Streams
- CDN Bandwidth Usage

---

## User Management

### User Operations

**View Users:**
```
Admin → Users → List All Users
```

**Filters:**
- Subscription tier
- Account status (active, suspended, cancelled)
- Registration date
- Last active

**User Actions:**
- View profile
- Edit details
- Suspend/Unsuspend account
- Reset password
- View activity log
- Manage subscriptions

### Bulk Operations

**Import Users:**
```bash
Admin → Users → Import → Upload CSV
```

CSV Format:
```
email,firstName,lastName,plan,status
user@example.com,John,Doe,premium,active
```

**Export Users:**
- Select filters
- Click "Export CSV"
- Download report

### Subscription Management

**Manual Subscription:**
1. Select user
2. Subscriptions → Add Subscription
3. Choose plan and duration
4. Apply discount (if any)
5. Set billing cycle
6. Confirm

**Refunds:**
1. Find subscription
2. Click "Refund"
3. Select refund type (full/partial)
4. Enter reason
5. Process refund

---

## Content Moderation

### Content Review Queue

**Access Queue:**
```
Admin → Content → Review Queue
```

**Filter by:**
- Flagged content
- New uploads (auto-review)
- Reported content
- Copyright claims

### Moderation Actions

**Approve Content:**
- Review details
- Check compliance
- Click "Approve"

**Reject Content:**
- Select reason
- Add notes
- Notify creator
- Click "Reject"

**Flag for Manual Review:**
- Add to escalation queue
- Assign to reviewer
- Set priority

### Automated Moderation

**AI Content Analysis:**
- Adult content detection
- Violence detection
- Copyright matching
- Inappropriate content flagging

**Configure Settings:**
```
Admin → Settings → Moderation
- Sensitivity levels (low, medium, high)
- Auto-actions (flag, block, approve)
- Notification settings
```

### Copyright Management

**DMCA Takedown Process:**
1. Receive DMCA notice
2. Verify claim validity
3. Remove infringing content
4. Notify uploader
5. Log action
6. Counter-claim process (if applicable)

---

## System Configuration

### Platform Settings

**General:**
- Platform name and branding
- Default language
- Time zone
- Date/time formats

**Features:**
- Enable/disable features
- Feature flags
- Beta features

**Limits:**
- Max file upload size
- Concurrent streams per account
- API rate limits
- Storage quotas

### Payment Configuration

**Payment Providers:**
- Stripe configuration
- PayPal setup
- Regional payment methods

**Subscription Plans:**
```
Admin → Settings → Plans
- Create/edit plans
- Set pricing
- Configure features
- Set trial periods
```

**Tax Settings:**
- Tax rates by region
- VAT configuration
- Tax exemptions

### Email Templates

**Manage Templates:**
```
Admin → Settings → Email Templates
```

Templates:
- Welcome email
- Password reset
- Subscription confirmation
- Payment failed
- Content published
- Moderation notices

### CDN Configuration

**CDN Settings:**
- Primary CDN provider
- Fallback CDNs
- Geographic routing
- Cache policies
- Purge cache

### DRM Configuration

**DRM Providers:**
- Widevine setup
- FairPlay configuration
- PlayReady settings

**License Server:**
- License server URL
- Certificate management
- Key rotation policy

---

## Analytics & Reporting

### Platform Analytics

**User Metrics:**
- Daily/Monthly Active Users (DAU/MAU)
- New registrations
- Churn rate
- Retention cohorts

**Content Metrics:**
- Total views
- Watch time
- Most viewed content
- Content library growth

**Revenue Metrics:**
- Revenue by source
- ARPU (Average Revenue Per User)
- MRR (Monthly Recurring Revenue)
- LTV (Lifetime Value)

**System Metrics:**
- Server uptime
- API response times
- Error rates
- CDN performance

### Custom Reports

**Create Report:**
1. Admin → Analytics → Reports
2. Select metrics
3. Choose date range
4. Apply filters
5. Generate report
6. Export (CSV, PDF, Excel)

**Schedule Reports:**
- Daily, weekly, monthly
- Email distribution
- Automated generation

---

## Security & Compliance

### Access Control

**Role Management:**
```
Admin → Security → Roles
```

**Default Roles:**
- Super Admin (full access)
- Admin (most features)
- Moderator (content review)
- Support (user assistance)
- Analyst (read-only)

**Create Custom Role:**
1. Define role name
2. Select permissions
3. Assign users
4. Save role

### Audit Logs

**View Logs:**
```
Admin → Security → Audit Logs
```

**Logged Actions:**
- User logins
- Admin operations
- Configuration changes
- Content moderation
- Payment transactions

**Log Retention:**
- 90 days (standard)
- 7 years (compliance mode)

### Security Settings

**Authentication:**
- Enforce MFA for admins
- Password policies
- Session timeout
- IP whitelisting

**API Security:**
- API key rotation
- Rate limiting
- IP restrictions
- Request signing

### Compliance

**GDPR Compliance:**
- Data export requests
- Right to be forgotten
- Consent management
- Data retention policies

**SOC 2 Controls:**
- Access controls
- Change management
- Incident response
- Backup verification

---

## Troubleshooting

### Common Issues

**Playback Failures:**
1. Check CDN status
2. Verify DRM licenses
3. Review error logs
4. Test from different regions

**Payment Issues:**
1. Check payment gateway status
2. Verify credentials
3. Review failed transactions
4. Contact payment provider

**Performance Issues:**
1. Check server metrics
2. Review database queries
3. Monitor CDN traffic
4. Scale resources if needed

### System Health Monitoring

**Dashboards:**
- Application health
- Database performance
- CDN status
- API status

**Alerts:**
- Configure alert thresholds
- Set notification channels
- On-call rotations

---

## Admin API Reference

### Authentication

```bash
curl -X POST https://admin-api.streamverse.io/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"***"}'
```

### User Management

**Get User:**
```bash
GET /v1/admin/users/{userId}
Authorization: Bearer {adminToken}
```

**Update User:**
```bash
PUT /v1/admin/users/{userId}
Content-Type: application/json

{
  "status": "suspended",
  "reason": "Terms violation"
}
```

### Content Management

**Get Content:**
```bash
GET /v1/admin/content/{contentId}
```

**Moderate Content:**
```bash
POST /v1/admin/content/{contentId}/moderate
{
  "action": "approve|reject|flag",
  "reason": "string",
  "notes": "string"
}
```

---

**Admin Guide Version**: 2.0
**© StreamVerse Inc. All rights reserved.**
