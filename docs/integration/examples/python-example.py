#!/usr/bin/env python3
"""VEDO EduTrack — Python client example (requests + JWT auth).

Usage:
    pip install requests
    python python-example.py

Shows the full flow: get a token, compute a route, check progress, and
receive a webhook subscription.
"""
import os
import sys

import requests

API_BASE = os.environ.get("API_BASE", "http://localhost:8080/api/v1")


def get_token() -> str:
    resp = requests.post(
        f"{API_BASE}/auth/token",
        json={"user_id": "user-1", "roles": ["learner"]},
        timeout=10,
    )
    resp.raise_for_status()
    return resp.json()["token"]


def main() -> int:
    token = get_token()
    headers = {"Authorization": f"Bearer {token}"}

    # 1. Compute a route.
    route = requests.post(
        f"{API_BASE}/routes/compute",
        headers=headers,
        json={"learner_id": "user-1", "goal_topic_id": "math-5-11"},
        timeout=10,
    )
    route.raise_for_status()
    print("Route steps:", len(route.json().get("steps", [])))

    # 2. Progress.
    progress = requests.get(
        f"{API_BASE}/learners/user-1/progress", headers=headers, timeout=10
    )
    progress.raise_for_status()
    print("Progress:", progress.json())

    # 3. Create a webhook subscription (https URL, 32+ char secret).
    secret = os.environ.get("WEBHOOK_SECRET", "unsafe")
    sub = requests.post(
        f"{API_BASE}/webhooks/subscriptions",
        headers=headers,
        json={
            "url": "https://your-app.example.com/hooks/edutrack",
            "event_types": ["module.mastered"],
            "secret": secret,
        },
        timeout=10,
    )
    if sub.status_code == 201:
        print("Webhook subscription created:", sub.json()["id"])
    else:
        print("Subscription error:", sub.status_code, sub.text)
        return 1

    return 0


if __name__ == "__main__":
    sys.exit(main())
