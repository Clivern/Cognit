# Cognit API

Cognit is a control plane for AI agents that speak the [A2A protocol](https://a2a-protocol.org/v1.0.0/specification/). A2A standardizes how you talk to an agent (Agent Card, messages, tasks). It does not standardize a registry. Cognit fills that gap: multi-tenant catalog, health, policy, and a gateway in front of A2A.

A2A v1.0 already has a `tenant` field (path parameter on operations, optional on `AgentInterface`). Tenancy lives in the URL, not only a header.

Base URL:

```
https://cognit.example.com
```

| Plane | Role |
| --- | --- |
| Control | Register, discover, health, policy |
| Data | Speak A2A to a resolved agent |

Auth on both planes: `Authorization: Bearer <token>`. Tenant is in the path. Optional `A2A-Version: 1.0` on data-plane calls.

---

## 1. Tenancy and identity

```
POST   /v1/tenants
GET    /v1/tenants
GET    /v1/tenants/{tenant}
PATCH  /v1/tenants/{tenant}
DELETE /v1/tenants/{tenant}

POST   /v1/tenants/{tenant}/api-keys
GET    /v1/tenants/{tenant}/api-keys
DELETE /v1/tenants/{tenant}/api-keys/{key_id}

GET    /v1/tenants/{tenant}/principals/me
```

A token is scoped to one tenant. Cross-tenant calls fail closed (`404`, not `403`) so tenant IDs stay private.

---

## 2. Control plane — catalog

Agents publish an A2A **Agent Card**. The registry stores the card plus runtime metadata (address, lease, health).

```
PUT    /v1/tenants/{tenant}/catalog/agents/{agent}
GET    /v1/tenants/{tenant}/catalog/agents/{agent}
DELETE /v1/tenants/{tenant}/catalog/agents/{agent}
GET    /v1/tenants/{tenant}/catalog/agents
       ?skill=invoice.extract
       &tag=finance
       &status=passing
       &capability=streaming

GET    /v1/tenants/{tenant}/catalog/agents/{agent}/card
PUT    /v1/tenants/{tenant}/catalog/agents/{agent}/card

GET    /v1/tenants/{tenant}/catalog/skills
GET    /v1/tenants/{tenant}/catalog/skills/{skill}
GET    /v1/tenants/{tenant}/catalog/skills/{skill}/agents
```

`PUT .../agents/{agent}` registers the agent. Body is the Agent Card plus instance info:

```json
{
  "card": {
    "name": "Invoice Extractor",
    "description": "Extracts line items from invoices",
    "version": "1.4.0",
    "capabilities": {
      "streaming": true,
      "pushNotifications": true,
      "extendedAgentCard": true
    },
    "defaultInputModes": ["text/plain", "application/pdf"],
    "defaultOutputModes": ["application/json"],
    "skills": [
      {
        "id": "invoice.extract",
        "name": "Extract invoice",
        "description": "Parse a PDF invoice into structured line items",
        "tags": ["finance", "ocr"]
      }
    ],
    "supportedInterfaces": [
      {
        "url": "https://invoice-worker.internal:8080/a2a/v1",
        "protocolBinding": "HTTP+JSON",
        "protocolVersion": "1.0",
        "tenant": "acme"
      }
    ]
  },
  "instance": {
    "id": "invoice-extractor-3",
    "address": "10.0.12.41",
    "port": 8080,
    "datacenter": "eu-west-1",
    "meta": { "model": "gpt-5", "region": "eu" }
  },
  "lease_ttl": "30s"
}
```

`GET .../catalog/agents` returns healthy instances and cards.

---

## 3. Control plane — self-register and leases

Workers talk to a local or cluster agent endpoint, not the full catalog UI.

```
PUT    /v1/tenants/{tenant}/agent/register
PUT    /v1/tenants/{tenant}/agent/deregister/{agent}

PUT    /v1/tenants/{tenant}/agent/instances/{instance}
DELETE /v1/tenants/{tenant}/agent/instances/{instance}

PUT    /v1/tenants/{tenant}/agent/instances/{instance}/renew
GET    /v1/tenants/{tenant}/agent/instances/{instance}

POST   /v1/tenants/{tenant}/sessions
GET    /v1/tenants/{tenant}/sessions/{session}
PUT    /v1/tenants/{tenant}/sessions/{session}/renew
DELETE /v1/tenants/{tenant}/sessions/{session}
```

Registration is lease-based. Miss the TTL and the instance is marked critical, then dropped from discovery.

---

## 4. Control plane — health

```
GET    /v1/tenants/{tenant}/health/agents
GET    /v1/tenants/{tenant}/health/agents/{agent}
GET    /v1/tenants/{tenant}/health/agents/{agent}/instances/{instance}

PUT    /v1/tenants/{tenant}/health/checks/{check}
DELETE /v1/tenants/{tenant}/health/checks/{check}
POST   /v1/tenants/{tenant}/health/checks/{check}/pass
POST   /v1/tenants/{tenant}/health/checks/{check}/warn
POST   /v1/tenants/{tenant}/health/checks/{check}/fail

GET    /v1/tenants/{tenant}/health/state/passing
GET    /v1/tenants/{tenant}/health/state/warning
GET    /v1/tenants/{tenant}/health/state/critical
```

Checks that matter for agents:

- **TTL heartbeat** — process is alive
- **A2A probe** — `GetTask` or a cheap `SendMessage` against a health skill
- **Card freshness** — `/.well-known/agent-card.json` still matches the registered card
- **Skill probe** — named skill still accepts the declared input modes

---

## 5. Control plane — discovery and saved queries

A2A’s own discovery story is a curated registry with an unspecified API. This is that API.

```
GET    /v1/tenants/{tenant}/discover
       ?skill=invoice.extract
       &tags=finance
       &input_mode=application/pdf
       &capability=streaming
       &passing=true
       &limit=5

POST   /v1/tenants/{tenant}/discover/resolve
```

Resolve body:

```json
{
  "skill": "invoice.extract",
  "strategy": "least_outstanding_tasks",
  "constraints": {
    "datacenter": "eu-west-1",
    "capabilities": ["streaming"]
  }
}
```

Response is the card the caller should use, with `supportedInterfaces[0].url` pointed at Cognit’s gateway (not the raw pod IP), plus the tenant key A2A expects:

```json
{
  "agent": "invoice-extractor",
  "instance": "invoice-extractor-3",
  "card": {
    "supportedInterfaces": [
      {
        "url": "https://cognit.example.com/v1/tenants/acme/a2a/agents/invoice-extractor",
        "protocolBinding": "HTTP+JSON",
        "protocolVersion": "1.0",
        "tenant": "acme"
      }
    ]
  }
}
```

Saved queries:

```
PUT    /v1/tenants/{tenant}/query/{name}
GET    /v1/tenants/{tenant}/query/{name}
DELETE /v1/tenants/{tenant}/query/{name}
GET    /v1/tenants/{tenant}/query/{name}/execute
```

---

## 6. Control plane — intentions

Who may call whom.

```
PUT    /v1/tenants/{tenant}/intentions
GET    /v1/tenants/{tenant}/intentions
GET    /v1/tenants/{tenant}/intentions/{id}
DELETE /v1/tenants/{tenant}/intentions/{id}

POST   /v1/tenants/{tenant}/intentions/check
```

```json
{
  "source": { "agent": "cfo-assistant" },
  "destination": { "agent": "invoice-extractor", "skill": "invoice.extract" },
  "action": "allow"
}
```

`POST .../intentions/check` is what the gateway calls on every A2A request. Default deny across tenants. Default deny across agents inside a tenant unless an intention exists.

---

## 7. Control plane — KV, watches, events

```
GET    /v1/tenants/{tenant}/kv/{key}
PUT    /v1/tenants/{tenant}/kv/{key}
DELETE /v1/tenants/{tenant}/kv/{key}
GET    /v1/tenants/{tenant}/kv?prefix=agents/invoice-extractor/

GET    /v1/tenants/{tenant}/watch/catalog?index={index}
GET    /v1/tenants/{tenant}/watch/health?index={index}
GET    /v1/tenants/{tenant}/watch/kv/{key}?index={index}

PUT    /v1/tenants/{tenant}/event/{name}
GET    /v1/tenants/{tenant}/event/{name}
```

Use KV for model pins, prompt refs, and feature flags — not secrets. Blocking `index` watches let agents long-poll and hot-reload cards and config.

---

## 8. Data plane — A2A well-known cards

A2A requires a card at `/.well-known/agent-card.json`. In a multi-tenant registry you host per-agent cards, and optionally a registry card.

```
GET  /v1/tenants/{tenant}/agents/{agent}/.well-known/agent-card.json
GET  /v1/tenants/{tenant}/agents/{agent}/extendedAgentCard

GET  /.well-known/agent-card.json
```

The last one is Cognit itself (directory / router agent). Per-agent cards are what callers cache. Serve `Content-Type: application/a2a+json` and `ETag` from `card.version`.

If you need a public hostname per agent:

```
GET  https://{agent}.{tenant}.cognit.example.com/.well-known/agent-card.json
```

That still resolves to the same card, with `supportedInterfaces` aimed at the gateway.

---

## 9. Data plane — A2A HTTP+JSON

Official A2A v1.0 REST bindings, namespaced by tenant and agent. `tenant` as a path param is what the spec already describes.

```
POST   /v1/tenants/{tenant}/a2a/agents/{agent}/message:send
POST   /v1/tenants/{tenant}/a2a/agents/{agent}/message:stream

GET    /v1/tenants/{tenant}/a2a/agents/{agent}/tasks
GET    /v1/tenants/{tenant}/a2a/agents/{agent}/tasks/{id}
POST   /v1/tenants/{tenant}/a2a/agents/{agent}/tasks/{id}:cancel
POST   /v1/tenants/{tenant}/a2a/agents/{agent}/tasks/{id}:subscribe

POST   /v1/tenants/{tenant}/a2a/agents/{agent}/tasks/{id}/pushNotificationConfigs
GET    /v1/tenants/{tenant}/a2a/agents/{agent}/tasks/{id}/pushNotificationConfigs
GET    /v1/tenants/{tenant}/a2a/agents/{agent}/tasks/{id}/pushNotificationConfigs/{configId}
DELETE /v1/tenants/{tenant}/a2a/agents/{agent}/tasks/{id}/pushNotificationConfigs/{configId}

GET    /v1/tenants/{tenant}/a2a/agents/{agent}/extendedAgentCard
```

JSON-RPC binding, same prefix, one POST:

```
POST   /v1/tenants/{tenant}/a2a/agents/{agent}
```

Methods: `SendMessage`, `SendStreamingMessage`, `GetTask`, `ListTasks`, `CancelTask`, `SubscribeToTask`, `CreateTaskPushNotificationConfig`, `GetTaskPushNotificationConfig`, `ListTaskPushNotificationConfigs`, `DeleteTaskPushNotificationConfig`, `GetExtendedAgentCard`.

The gateway:

1. Authenticates the caller
2. Checks intentions
3. Picks a passing instance (or the instance sticky to `task.id`)
4. Forwards the A2A call
5. Records health and outstanding-task load

---

## 10. Skill-addressed A2A

Callers often do not know the agent name. They know a skill.

```
POST   /v1/tenants/{tenant}/a2a/skills/{skill}/message:send
POST   /v1/tenants/{tenant}/a2a/skills/{skill}/message:stream
GET    /v1/tenants/{tenant}/a2a/skills/{skill}/tasks/{id}
POST   /v1/tenants/{tenant}/a2a/skills/{skill}/tasks/{id}:cancel
POST   /v1/tenants/{tenant}/a2a/skills/{skill}/tasks/{id}:subscribe
```

`{skill}` is `invoice.extract`. Cognit resolves, checks intentions, picks a healthy instance, and forwards the A2A call. Task IDs stay sticky to the instance that accepted the first message.

---

## 11. Cluster and ops

```
GET    /v1/status/leader
GET    /v1/status/peers
GET    /v1/health

GET    /v1/tenants/{tenant}/datacenters
GET    /v1/tenants/{tenant}/nodes
GET    /v1/tenants/{tenant}/metrics
```

`GET /v1/health` is process liveness and is not tenant-scoped.

---

## Request shape

Every tenant-scoped call:

```
Authorization: Bearer <tenant-scoped token>
A2A-Version: 1.0
Content-Type: application/json
```

A2A send through the gateway:

```http
POST /v1/tenants/acme/a2a/agents/invoice-extractor/message:send HTTP/1.1
Authorization: Bearer ...
A2A-Version: 1.0
Content-Type: application/json

{
  "message": {
    "role": "ROLE_USER",
    "parts": [
      { "text": "Extract line items" }
    ]
  }
}
```

---

## Build order

Thin slice that is already a useful A2A core:

1. `PUT/GET/DELETE /v1/tenants/{tenant}/catalog/agents/{agent}`
2. `PUT .../agent/instances/{instance}/renew` (TTL)
3. `GET .../discover?skill=`
4. `GET .../agents/{agent}/.well-known/agent-card.json`
5. `POST .../a2a/agents/{agent}/message:send` and `GET .../tasks/{id}`
6. Intentions `check` on the gateway

KV, watches, saved queries, and skill-addressed routing come after that.

A2A’s job is the card plus message and task verbs. Cognit’s job is who is registered, who is healthy, who is allowed to call, and which instance gets the RPC — with tenant as a hard path boundary.
