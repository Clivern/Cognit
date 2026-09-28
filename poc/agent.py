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
    }
    return replies[skill_id]
