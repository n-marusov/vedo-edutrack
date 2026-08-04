# Webhook Guide — Receiving EduTrack Events

EduTrack sends webhook events when key domain milestones occur:
`module.mastered`, `plan.deviated`, `route.recalculated`, and
`standard.risk_detected`.

## Subscribing

```
POST /api/v1/webhooks/subscriptions
Content-Type: application/json
Authorization: Bearer <token>

{
  "url": "https://your-app.example.com/hooks/edutrack",
  "event_types": ["module.mastered"],
  "secret": "<signing-secret-32-chars-min>"
}
```

Response (201):

```json
{
  "id": "sub-uuid",
  "url": "https://your-app.example.com/hooks/edutrack",
  "event_types": ["module.mastered"],
  "active": true,
  "secret": "<your-secret>"
}
```

## Delivery format

Every event is delivered as an HTTP POST with a JSON body:

```json
{
  "event_id": "d4c1e2a5-1f0e-4c3a-9b2f-6a3c8e7d1a2b",
  "event_type": "module.mastered",
  "timestamp": 1694123456,
  "data": {
    "learner_id": "user-1",
    "module_id": "percent",
    "mastered_at": "2026-09-05T12:00:00Z"
  }
}
```

## Signature verification

Each request carries an `X-Vedo-Signature` header:

```
X-Vedo-Signature: t=1694123456,v1=<hex-signature>
```

To verify:

1. Extract `t` (unix timestamp) and `v1` (hex HMAC-SHA256) from the header.
2. Compute: `hex(HMAC-SHA256(secret, "t.<payload>"))`
3. Compare with `v1`. Use a constant-time comparison.

Example (Node.js):

```js
const crypto = require('crypto');

function verifySignature(header, secret, payload) {
  const parts = header.split(',');
  const t = parts.find(p => p.startsWith('t=')).slice(2);
  const sig = parts.find(p => p.startsWith('v1=')).slice(3);
  const expected = crypto
    .createHmac('sha256', secret)
    .update(`${t}.${payload}`)
    .digest('hex');
  return crypto.timingSafeEqual(Buffer.from(sig), Buffer.from(expected));
}
```

## Idempotency

Duplicate events (same `event_id`) are delivered only once. If you receive a
duplicate, ignore it.

## Retry policy

| Attempt | Delay (exponential backoff) |
|---------|----------------------------|
| 1       | 1s                         |
| 2       | 2s                         |
| 3       | 4s                         |
| 4       | 8s                         |
| 5       | 16s                        |

After 5 consecutive failures the subscription is deactivated.

## Testing with curl

```bash
# Create a ping event
curl -X POST http://localhost:8080/api/v1/webhooks/subscriptions/<id>/ping \
  -H "Authorization: Bearer $TOKEN"

# View delivery history
curl http://localhost:8080/api/v1/webhooks/subscriptions/<id>/deliveries \
  -H "Authorization: Bearer $TOKEN" | jq
```

## See also

- `examples/webhook-receiver-example.js` — runnable Node receiver
- `quickstart.md` — get started in 5 minutes
