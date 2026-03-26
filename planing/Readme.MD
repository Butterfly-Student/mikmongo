# Planning Revisi — ISP AI Platform
> Berdasarkan kondisi aktual: Go billing sudah production

---

## Status Saat Ini

```
✅ SELESAI — Go Backend ISP Billing
   ├── Billing & invoice otomatis
   ├── Isolir & reaktivasi customer
   ├── PPPoE management ke MikroTik
   ├── WhatsApp notifikasi
   └── Hotspot Mikhmon-style (data di DB, scripting di MikroTik)

🔲 BELUM — GenieACS + OLT integration (tambahan di Go)
🔲 BELUM — Python AI Agent (TCP stream + LangGraph + WA interaksi)
🔲 BELUM — TypeScript Dashboard
```

---

## Strategi: Paralel 3 Track

Karena Go sudah production, 3 track berikut bisa dikerjakan **bersamaan** tanpa saling blokir:

```
Track A  ──── GenieACS + OLT ──────────────────────── (2–3 minggu)
Track B  ──────── Python Agent ─────────────────────────────────── (5–6 minggu)
Track C  ──────── TypeScript Dashboard ──────────────────────────── (4–5 minggu)
         W1       W2       W3       W4       W5       W6
```

Jika dikerjakan 1 developer: kerjakan Track A dulu (cepat, hasilnya langsung dipakai Python), lalu B dan C bisa overlap.

---

## Track A — GenieACS + OLT Integration (Go)
**Estimasi: 2–3 minggu | Prioritas: Tinggi**

> Go sudah punya fondasi MikroTik. Track ini menambah 2 integrasi baru ke layer yang sama.

### A.1 GenieACS NBI Integration
GenieACS sudah jalan sebagai TR-069 server. Go cukup bicara ke GenieACS NBI (REST API).

```
Endpoint GenieACS NBI yang dipakai:
GET  /devices                    → list semua ONU terdaftar
GET  /devices/{id}               → detail device: serial, firmware, status
POST /devices/{id}/tasks         → trigger: inform, reboot, factory reset
GET  /devices?query={filter}     → filter by OUI, product class
```

Checklist:
- [ ] HTTP client ke GenieACS NBI (internal, tidak expose ke publik)
- [ ] Get device info by serial number
- [ ] Get last inform time & connection status
- [ ] Trigger inform/reboot ONU
- [ ] Bulk status: semua ONU yang offline > N menit
- [ ] Provisioning ONU baru: set parameter TR-069 via task
- [ ] Link ONU serial ke customer di DB

### A.2 OLT Integration
OLT bisa diakses via SNMP (read) dan SSH/Telnet (read + write).

```
SNMP (read-only, polling):
- ONU registration status (ifOperStatus)
- Optical receive power per ONU (GPON MIB)
- OLT uptime, port status

SSH/Telnet (command):
- show onu optical-info (optical power)
- show onu state (status per port)
- interface onu vlan (konfigurasi VLAN)
```

Checklist:
- [ ] SNMP client (gosnmp): get optical power per ONU
- [ ] SSH client: connect ke OLT, eksekusi show commands
- [ ] Parser output CLI OLT (per vendor: ZTE, Huawei, Fiberhome)
- [ ] Simpan optical power history ke DB (time-series per ONU)
- [ ] Alert threshold: rx_power < -27 dBm → insert ke security_alerts

### A.3 API Endpoint Baru (untuk Python & TypeScript)
```
# Internal (untuk Python agent)
GET  /internal/genieacs/device/{serial}      → detail ONU
GET  /internal/genieacs/offline              → list ONU offline
POST /internal/genieacs/reboot/{serial}      → reboot ONU
GET  /internal/olt/optical-power/{serial}    → optical power realtime
GET  /internal/olt/onu-status                → semua ONU + status

# Public (untuk TypeScript dashboard)
GET  /api/onu/list                           → list ONU + customer + signal
GET  /api/onu/{id}/history                   → history optical power
GET  /api/olt/status                         → status semua OLT
```

**Deliverable Track A:** Python agent bisa query GenieACS & OLT via Go API. Dashboard bisa tampilkan status ONU dan optical power.

---

## Track B — Python AI Agent
**Estimasi: 5–6 minggu | Prioritas: Tinggi**

> Ini bagian yang paling baru. Panduan lengkap dari nol.

### Struktur Project Python

```
python-agent/
├── main.py                  ← entry point (FastAPI + Celery)
├── config.py                ← env vars (Go API URL, LLM key, dll)
│
├── mikrotik/
│   ├── stream_manager.py    ← MultiRouterStreamManager (asyncio)
│   ├── listeners.py         ← 4 stream per router
│   └── anomaly_detector.py  ← event handler + detection logic
│
├── agent/
│   ├── graph.py             ← LangGraph agent graph definition
│   ├── tools.py             ← semua tool definitions
│   ├── router.py            ← classify_intent (rule-based vs LLM)
│   └── memory.py            ← context per sesi (Redis)
│
├── workers/
│   ├── celery_app.py        ← Celery init + Redis broker
│   ├── worker_wa_operator.py
│   ├── worker_wa_customer.py
│   ├── worker_stream.py     ← jalankan stream_manager
│   └── worker_security.py   ← proses security events
│
├── api/
│   ├── routes.py            ← /agent/chat, /agent/target, /agent/status
│   └── auth.py              ← validasi API key dari Go
│
└── rag/
    ├── setup.py             ← ingest docs ke ChromaDB
    ├── retriever.py         ← query RAG
    └── docs/                ← MikroTik wiki, ISP runbook
```

### B.1 Mulai Dari Mana? — Urutan yang Disarankan

```
Minggu 1: Setup + Go API client + TCP stream dasar
Minggu 2: Anomaly detection (brute force, traffic spike)
Minggu 3: LangGraph agent + tools dasar
Minggu 4: WhatsApp interaksi (operator + customer)
Minggu 5: RAG setup + Python Agent API untuk TypeScript
Minggu 6: Celery workers + integrasi lengkap + test
```

### B.2 Langkah 1 — Setup & Go API Client

```python
# config.py
import os

GO_API_URL = os.getenv("GO_API_URL", "http://localhost:8080")
GO_API_KEY = os.getenv("GO_API_KEY", "")
LLM_API_KEY = os.getenv("ANTHROPIC_API_KEY", "")
REDIS_URL = os.getenv("REDIS_URL", "redis://localhost:6379")
MIKROTIK_DEFAULT = os.getenv("MIKROTIK_DEFAULT_ID", "router-01")
```

```python
# agent/tools.py — Go API client
import httpx
from config import GO_API_URL, GO_API_KEY

async def go_get(path: str) -> dict:
    async with httpx.AsyncClient() as client:
        resp = await client.get(
            f"{GO_API_URL}{path}",
            headers={"X-API-Key": GO_API_KEY},
            timeout=10.0
        )
        resp.raise_for_status()
        return resp.json()

async def go_post(path: str, data: dict) -> dict:
    async with httpx.AsyncClient() as client:
        resp = await client.post(
            f"{GO_API_URL}{path}",
            json=data,
            headers={"X-API-Key": GO_API_KEY},
            timeout=10.0
        )
        resp.raise_for_status()
        return resp.json()
```

### B.3 Langkah 2 — MikroTik TCP Stream

```python
# Install: pip install librouteros

# mikrotik/stream_manager.py
import asyncio
import librouteros

class MikroTikStreamListener:
    def __init__(self, router: dict):
        self.router = router
        self.conn = None

    async def connect(self):
        self.conn = librouteros.connect(
            host=self.router['ip'],
            username=self.router['username'],
            password=self.router['password'],
            port=8728
        )

    async def listen_log(self, callback):
        """Push setiap ada log entry baru"""
        for entry in self.conn('/log/listen'):
            await callback(entry, self.router['id'])

    async def listen_pppoe(self, callback):
        """Push setiap ada perubahan PPPoE session"""
        for event in self.conn('/ppp/active/listen'):
            await callback(event, self.router['id'])

    async def listen_traffic(self, iface: str, callback):
        """Push traffic stats setiap 1 detik"""
        for data in self.conn(
            '/interface/monitor-traffic',
            **{'interface': iface, 'once': 'false'}
        ):
            await callback(data, self.router['id'])


class MultiRouterStreamManager:
    def __init__(self, detector, go_client):
        self.detector = detector
        self.go = go_client
        self.listeners = {}
        self.current_target = None  # in-memory, per sesi

    async def startup(self):
        """Ambil registry dari Go, start semua stream"""
        routers = await go_get('/internal/mikrotik/list')
        await asyncio.gather(*[
            self._start_with_reconnect(r) for r in routers
        ])

    async def _start_with_reconnect(self, router: dict):
        backoff = 1
        while True:
            try:
                listener = MikroTikStreamListener(router)
                await listener.connect()
                self.listeners[router['id']] = listener
                # 3 stream paralel per router
                await asyncio.gather(
                    listener.listen_log(self.detector.on_log),
                    listener.listen_pppoe(self.detector.on_pppoe),
                    listener.listen_traffic('ether1', self.detector.on_traffic),
                    return_exceptions=True
                )
                backoff = 1
            except Exception as e:
                await go_post('/internal/router/connection-lost',
                              {'router_id': router['id'], 'error': str(e)})
                await asyncio.sleep(min(backoff, 60))
                backoff *= 2

    def switch_target(self, router_id: str):
        """Switch current target — in-memory only, tidak ubah DB"""
        if router_id in self.listeners:
            self.current_target = router_id
            return True
        return False
```

### B.4 Langkah 3 — Anomaly Detector

```python
# mikrotik/anomaly_detector.py

class RealtimeAnomalyDetector:
    def __init__(self, go_client):
        self.go = go_client
        self.login_fails = {}   # IP → [timestamps]
        self.pppoe_events = {}  # user → [timestamps]

    async def on_log(self, entry: dict, router_id: str):
        msg = entry.get('message', '').lower()
        if 'login failure' in msg:
            ip = self._extract_ip(msg)
            if ip:
                await self._track_brute_force(ip, router_id)
        if 'link down' in msg:
            await self._alert('interface_down', router_id,
                              {'message': msg})

    async def on_pppoe(self, event: dict, router_id: str):
        if event.get('.type') == 'remove':
            user = event.get('name', '')
            await self._track_flapping(user, router_id)

    async def on_traffic(self, data: dict, router_id: str):
        rx = int(data.get('rx-bits-per-second', 0))
        # Threshold bisa disimpan per router di Redis
        if rx > 800_000_000:  # > 800 Mbps = anomali
            await self._alert('traffic_spike', router_id,
                              {'rx_bps': rx})

    async def _track_brute_force(self, ip: str, router_id: str):
        key = f"{router_id}:{ip}"
        now = asyncio.get_event_loop().time()
        self.login_fails.setdefault(key, [])
        # Sliding window 60 detik
        self.login_fails[key] = [
            t for t in self.login_fails[key] if now - t < 60
        ] + [now]
        count = len(self.login_fails[key])
        if count == 10 or count == 30:  # alert di 10 dan 30
            await self._alert('brute_force', router_id,
                              {'ip': ip, 'count': count,
                               'severity': 'high' if count >= 30 else 'medium'})

    async def _alert(self, type_: str, router_id: str, extra: dict):
        await go_post('/internal/security/alert',
                      {'type': type_, 'router_id': router_id, **extra})
```

### B.5 Langkah 4 — LangGraph Agent

```python
# Install: pip install langchain-core langchain-anthropic langgraph

# agent/graph.py
from langgraph.graph import StateGraph, END
from langchain_anthropic import ChatAnthropic
from agent.tools import (
    get_customer_info, get_billing, get_pppoe_status,
    get_onu_status, get_optical_power, get_alerts,
    mikrotik_query, isolate_customer
)

# Tools yang tersedia untuk agent
TOOLS = [
    get_customer_info,  # → go GET /internal/customers/{id}
    get_billing,        # → go GET /api/billing
    get_pppoe_status,   # → go GET /api/mikrotik/sessions
    get_onu_status,     # → go GET /internal/genieacs/device/{serial}
    get_optical_power,  # → go GET /internal/olt/optical-power/{serial}
    get_alerts,         # → go GET /api/alerts
    mikrotik_query,     # → librouteros read-only query
]

OPERATOR_TOOLS = TOOLS + [
    isolate_customer,   # → go POST /internal/billing/isolate/{id}
    # tambah tool eksekusi lain untuk operator
]

llm = ChatAnthropic(model="claude-sonnet-4-20250514")

def build_agent_graph(mode: str = "operator"):
    tools = OPERATOR_TOOLS if mode == "operator" else TOOLS
    llm_with_tools = llm.bind_tools(tools)

    def agent_node(state):
        messages = state["messages"]
        response = llm_with_tools.invoke(messages)
        return {"messages": [response]}

    graph = StateGraph(dict)
    graph.add_node("agent", agent_node)
    # tambah tool node, dll
    graph.set_entry_point("agent")
    return graph.compile()
```

### B.6 Langkah 5 — Intent Router (LLM hanya jika perlu)

```python
# agent/router.py
# Pertanyaan sederhana tidak perlu ke LLM → hemat cost

SIMPLE_PATTERNS = {
    "tagihan": ["tagihan", "bayar", "invoice", "nunggak"],
    "status":  ["internet", "mati", "lambat", "gangguan"],
    "isolir":  ["terisolir", "diblokir", "tidak bisa"],
}

def classify_intent(message: str) -> str:
    msg = message.lower()
    for intent, keywords in SIMPLE_PATTERNS.items():
        if any(kw in msg for kw in keywords):
            return intent
    return "llm"  # fallback ke LLM

async def handle_customer_message(message: str, customer_id: str) -> str:
    intent = classify_intent(message)

    if intent == "tagihan":
        # Rule-based: tidak ke LLM
        billing = await go_get(f"/api/billing/{customer_id}/latest")
        return format_billing_response(billing)

    elif intent == "status":
        # Cek PPPoE dan ONU, jawab tanpa LLM
        session = await go_get(f"/api/mikrotik/session/{customer_id}")
        return format_status_response(session)

    else:
        # Hanya pertanyaan kompleks ke LLM
        return await run_agent(message, customer_id, mode="customer")
```

### B.7 Langkah 6 — WhatsApp Handler

```python
# workers/worker_wa_operator.py

@celery_app.task
async def handle_operator_message(phone: str, message: str):
    # Operator punya akses penuh
    context = await get_operator_context(phone)
    intent = classify_operator_intent(message)

    if intent == "switch_router":
        # Switch target Python agent
        router_name = extract_router_name(message)
        success = stream_manager.switch_target(router_name)
        return f"{'✅' if success else '❌'} Target: {router_name}"

    elif intent == "query":
        # Tanya agent dengan akses penuh
        response = await run_agent(message, context, mode="operator")
        await send_wa(phone, response)

# workers/worker_wa_customer.py

@celery_app.task
async def handle_customer_message(phone: str, message: str):
    # Validasi customer via Go
    customer = await go_get(f"/internal/customers/by-phone/{phone}")
    if not customer:
        await send_wa(phone, "Nomor tidak terdaftar.")
        return

    response = await handle_customer_intent(message, customer['id'])
    await send_wa(phone, response)
```

### B.8 Langkah 7 — Python Agent API (untuk TypeScript)

```python
# api/routes.py
from fastapi import FastAPI, Header, HTTPException

app = FastAPI()

@app.post("/agent/chat")
async def chat(body: ChatRequest, x_api_key: str = Header(...)):
    validate_key(x_api_key)
    response = await run_agent(body.message, mode=body.role)
    return {"answer": response.text, "data": response.structured_data}

@app.post("/agent/target")
async def set_target(body: TargetRequest, x_api_key: str = Header(...)):
    validate_key(x_api_key)
    success = stream_manager.switch_target(body.router_id)
    if not success:
        raise HTTPException(404, "Router not found")
    return {"status": "ok", "current_target": body.router_id}

@app.get("/agent/status")
async def get_status():
    return {
        "current_target": stream_manager.current_target,
        "available_routers": list(stream_manager.listeners.keys()),
        "stream_health": stream_manager.get_health()
    }

@app.get("/agent/stream-health")
async def stream_health():
    return stream_manager.get_health()
    # return: {router_id: {connected, last_event, streams: [...]}}
```

**Deliverable Track B:** Realtime monitoring aktif. Brute force < 5 detik. Agent via WhatsApp untuk operator & customer. API untuk TypeScript siap.

---

## Track C — TypeScript Dashboard
**Estimasi: 4–5 minggu | Prioritas: Paralel dengan B**

> Track C bisa mulai sejak Track A selesai karena Go API sudah ada.

### C.1 Setup

```bash
npx create-next-app@latest ts-dashboard --typescript --tailwind --app
cd ts-dashboard
npx shadcn-ui@latest init
npm install @tanstack/react-query @tanstack/react-virtual recharts socket.io-client
```

### C.2 Urutan Halaman

```
Minggu 1: Layout + Auth + Overview dashboard (data dari Go)
Minggu 2: Customer list (virtual scroll) + Billing UI
Minggu 3: Network monitor + ONU/optical power
Minggu 4: Agent chat panel + switch target
Minggu 5: WebSocket realtime + alert feed + polish
```

### C.3 Halaman yang Dibangun

```
/dashboard          → overview: stats, alert feed realtime, session count
/customers          → list 1500+ dengan virtual scroll + filter + pagination
/customers/[id]     → detail: info, tagihan, PPPoE session, ONU status
/billing            → invoice list, laporan pendapatan, history
/network            → semua MikroTik + label "Agent monitoring: X"
/network/[id]       → detail router: traffic, PPPoE sessions, log events
/onu                → list ONU + optical power + color coding signal
/alerts             → history anomali + security events
/agent              → chat panel ke Python agent (POST /agent/chat)
/settings           → MikroTik registry, threshold, template WA
```

### C.4 Chat Panel ke Python Agent

```typescript
// hooks/useAgentChat.ts
import { useState } from 'react'

interface Message {
  role: 'operator' | 'agent'
  content: string
  data?: any  // structured data dari agent response
}

export function useAgentChat() {
  const [messages, setMessages] = useState<Message[]>([])
  const [loading, setLoading] = useState(false)

  const sendMessage = async (text: string) => {
    setLoading(true)
    setMessages(prev => [...prev, { role: 'operator', content: text }])

    // POST langsung ke Python agent — bukan lewat Go
    const res = await fetch(`${process.env.NEXT_PUBLIC_AGENT_URL}/agent/chat`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json',
                 'X-API-Key': process.env.NEXT_PUBLIC_AGENT_KEY! },
      body: JSON.stringify({ message: text, role: 'operator' })
    })
    const data = await res.json()
    setMessages(prev => [...prev, {
      role: 'agent',
      content: data.answer,
      data: data.data  // untuk render grafik/tabel
    }])
    setLoading(false)
  }

  return { messages, sendMessage, loading }
}
```

### C.5 Switch Target Agent dari Dashboard

```typescript
// components/NetworkMonitor.tsx
async function setAgentTarget(routerId: string) {
  // Panggil Python agent API langsung — tidak lewat Go
  const res = await fetch(`${AGENT_URL}/agent/target`, {
    method: 'POST',
    headers: { 'X-API-Key': AGENT_KEY, 'Content-Type': 'application/json' },
    body: JSON.stringify({ router_id: routerId })
  })
  const data = await res.json()
  // Update label di UI
  setCurrentAgentTarget(data.current_target)
}
```

### C.6 WebSocket Realtime Alert

```typescript
// hooks/useAlertStream.ts
import { useEffect } from 'react'
import { io } from 'socket.io-client'

export function useAlertStream(onAlert: (alert: any) => void) {
  useEffect(() => {
    // Connect ke WebSocket Go backend
    const socket = io(process.env.NEXT_PUBLIC_GO_WS_URL!)
    socket.on('alert', onAlert)
    return () => { socket.disconnect() }
  }, [])
}
```

**Deliverable Track C:** Full dashboard. 1500 customer load < 3 detik. Chat panel agent aktif. Alert feed realtime.

---

## Timeline Paralel (Revisi)

```
              W1    W2    W3    W4    W5    W6
Track A (Go)  ████  ████  ██
Track B (Py)        ████  ████  ████  ████  ██
Track C (TS)  ████  ████  ████  ████  ██

Catatan:
- Track A mulai W1 karena paling cepat selesai (2–3 minggu)
- Track B mulai W1–2 karena independent dari A (bisa pakai mock data dulu)
- Track C mulai W1 karena Go API sudah ada untuk data billing & customer
- Track B dan C saling terhubung di W4–5 saat Python API sudah ready
```

---

## Dependensi Antar Track

```
Track A (Go) → Track B (Python)
  Python butuh endpoint baru: /internal/genieacs/*, /internal/olt/*
  Tapi Python bisa development dengan Go API yang sudah ada dulu

Track B (Python) → Track C (TypeScript)
  Dashboard butuh /agent/status, /agent/chat, /agent/target
  TypeScript bisa development dengan mock data sampai Python API ready

Track A (Go) → Track C (TypeScript)
  Dashboard halaman ONU butuh /api/onu/list dari Track A
  Halaman billing, customer, alerts sudah ada dari Go yang sekarang
```

---

## Checklist Immediate Next Steps

### Minggu Ini (Mulai Sekarang):

**Track A:**
- [ ] Buat HTTP client untuk GenieACS NBI di Go
- [ ] Test koneksi ke GenieACS: GET /devices
- [ ] Buat model ONU di DB (onu_devices, link ke customer)

**Track B:**
- [ ] `pip install librouteros langchain-core langchain-anthropic langgraph fastapi celery redis httpx`
- [ ] Setup project structure sesuai di atas
- [ ] Buat `config.py` dengan env vars
- [ ] Test koneksi TCP stream ke 1 MikroTik dengan `librouteros`
- [ ] Buat `go_get()` / `go_post()` client dan test ke Go API yang sudah ada

**Track C:**
- [ ] `npx create-next-app@latest ts-dashboard`
- [ ] Setup shadcn + react-query + react-virtual
- [ ] Buat halaman `/customers` dengan data dari Go API yang sudah ada
- [ ] Buat halaman `/billing` — data sudah tersedia dari Go