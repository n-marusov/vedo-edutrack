#!/usr/bin/env node
/**
 * VEDO EduTrack — Webhook receiver example (Node.js).
 *
 * Start with:
 *   npm install express
 *   node webhook-receiver-example.js
 *
 * Then point your webhook subscription to http://localhost:9099/hooks/edutrack.
 */
const crypto = require("crypto");
const express = require("express");

const PORT = 9099;
const SECRET = process.env.WEBHOOK_SECRET || "01234567890123456789012345678901";

const app = express();

// The raw body is needed for HMAC verification, but for simplicity we parse
// JSON and verify against the stringified body. In production, use
// raw-body middleware.
app.post("/hooks/edutrack", express.json(), (req, res) => {
  const signature = req.headers["x-vedo-signature"];
  if (!signature) {
    console.error("Missing X-Vedo-Signature");
    return res.status(400).send("Missing signature");
  }

  // Verify HMAC
  const parts = signature.split(",");
  const t = parts.find((p) => p.startsWith("t="))?.slice(2);
  const sig = parts.find((p) => p.startsWith("v1="))?.slice(3);

  if (!t || !sig) {
    return res.status(400).send("Invalid signature format");
  }

  const payload = JSON.stringify(req.body);
  const expected = crypto
    .createHmac("sha256", SECRET)
    .update(`${t}.${payload}`)
    .digest("hex");

  if (!crypto.timingSafeEqual(Buffer.from(sig), Buffer.from(expected))) {
    console.error("Signature mismatch");
    return res.status(403).send("Invalid signature");
  }

  // Process event
  console.log("Received:", req.body.event_type, req.body.event_id);
  console.log("  Data:", JSON.stringify(req.body.data));

  // Acknowledge (200 = success, 500+ = retry)
  res.status(200).send("ok");
});

app.listen(PORT, () => {
  console.log(`Webhook receiver listening on :${PORT}`);
  console.log(`Secret: ${SECRET}`);
});
