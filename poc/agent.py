"""Fake replies for every skill. uv run uvicorn agent:app --port 9100"""

import json
from pathlib import Path

from fastapi import FastAPI

SKILLS = json.loads(Path(__file__).with_name("skills.json").read_text())

app = FastAPI()


def fake(skill_id, args):
    user = args.get("user_id") or "u_123"
    email = args.get("email") or "jane@acme.io"
    team = args.get("team_id") or "sre-02"
    replies = {
        "user.profile.get": {
            "user_id": user,
            "name": "Jane",
            "email": "jane@acme.io",
            "role": "member",
        },
        "user.profile.update": {
            "user_id": user,
            "name": args.get("name") or "Jane",
            "email": args.get("email") or "jane@acme.io",
            "updated": True,
        },
        "user.password.reset": {"user_id": user, "sent": True},
        "user.session.list": {
            "user_id": user,
            "sessions": [{"session_id": "sess_1", "device": "laptop"}],
        },
        "user.session.revoke": {
            "user_id": user,
            "session_id": args.get("session_id") or "sess_1",
            "revoked": True,
        },
        "team.members.list": {
            "team_id": team,
            "members": ["jane@acme.io", "joe@gmail.com"],
        },
        "team.invite.send": {"email": email, "team_id": team, "invited": True},
        "team.member.remove": {"email": email, "team_id": team, "removed": True},
        "billing.renewal.get": {
            "user_id": user,
            "plan": "team",
            "renews_on": "2026-11-01",
        },
        "billing.plan.change": {
            "user_id": user,
            "plan": args.get("plan") or "team",
            "changed": True,
        },
        "billing.plan.cancel": {"user_id": user, "cancelled": True},
        "billing.invoice.list": {
            "user_id": user,
            "invoices": [
                {"invoice_id": "in_1", "amount": "20.00"},
                {"invoice_id": "in_2", "amount": "20.00"},
            ],
        },
        "billing.invoice.get": {
            "invoice_id": args.get("invoice_id") or "in_1",
            "amount": "20.00",
            "status": "paid",
        },
        "billing.payment_method.update": {"user_id": user, "card": "visa 4242"},
        "billing.refund.create": {
            "user_id": user,
            "invoice_id": args.get("invoice_id") or "in_1",
            "refunded": True,
        },
        "usage.summary.get": {"user_id": user, "used": 1200, "included": 5000},
        "support.ticket.create": {
            "user_id": user,
            "ticket_id": "t_1",
            "text": args.get("text") or "",
            "status": "open",
        },
        "support.ticket.list": {
            "user_id": user,
            "tickets": [{"ticket_id": "t_1", "status": "open"}],
        },
        "support.ticket.status": {
            "ticket_id": args.get("ticket_id") or "t_1",
            "status": "open",
        },
        "product.status.get": {"up": True},
        "product.outage.list": {"outages": []},
        "notification.prefs.get": {
            "user_id": user,
            "emails": ["billing", "product"],
        },
        "notification.prefs.update": {
            "user_id": user,
            "emails": args.get("emails") or [],
            "updated": True,
        },
        "audit.log.list": {
            "user_id": user,
            "events": [{"action": "login", "at": "2026-09-28T10:00:00Z"}],
        },
    }
    return replies[skill_id]


@app.get("/.well-known/agent-card.json")
def agent_card():
    return {
        "name": "Agent",
        "description": "Fake replies for every skill",
        "version": "1.0.0",
        "protocolVersion": "1.0",
        "url": "http://127.0.0.1:9100/a2a/v1",
        "skills": [
            {"id": skill["id"], "description": skill["description"]} for skill in SKILLS
        ],
    }


@app.post("/a2a/v1/message:send")
def message_send(body: dict):
    skill_id = body["skill"]
    data = body["message"]["parts"][0].get("data") or {}
    return {"message": {"role": "ROLE_AGENT", "parts": [{"data": fake(skill_id, data)}]}}
