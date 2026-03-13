---
name: mikrotik-routeros-cli
description: Complete Mikrotik RouterOS CLI reference for ISP and RT/RW Net — queue, interface (ethernet/bridge/vlan/pppoe), ppp (secret/profile/active), dns, ip pool/address/route/firewall/hotspot, system logging/scheduler/script, and tools. Use when writing RouterOS config commands, scripts, or automation code for Mikrotik devices.
dependencies: mikrotik-routeros>=6.x
---

# Mikrotik RouterOS CLI Reference

## Overview

This Skill provides exact RouterOS CLI command reference derived from a live RB750G device (RouterOS 6.49.11, platform G-Net, MIPS arch). Apply whenever the task involves:
- Writing or generating RouterOS configuration scripts
- Monitoring interfaces, queues, PPP sessions, hotspot users
- Managing IP addressing, routing, firewall (filter/nat/mangle/raw), address-lists
- Configuring hotspot captive portal, PPP/PPPoE secrets and profiles
- Reading system resources, logs, scheduler scripts
- Troubleshooting with ping, torch, traceroute, profile tools
- Building Go programs that automate RouterOS via API

> **When accessing RouterOS programmatically from Go**, always read the **go-routeros** Skill alongside this one. That Skill covers the `go-routeros/routeros/v3` library, connection handling, `Run` vs `Listen`, async mode, and CLI-to-API path mapping. The two Skills complement each other: this Skill defines *what* commands/parameters exist; the go-routeros Skill defines *how* to send them from Go code.

### Official Mikrotik Documentation References
- RouterOS API protocol: https://help.mikrotik.com/docs/spaces/ROS/pages/47579160/API
- Queues: https://help.mikrotik.com/docs/spaces/ROS/pages/328088/Queues
- PPP/AAA: https://help.mikrotik.com/docs/spaces/ROS/pages/132350049/PPP+AAA
- HotSpot: https://help.mikrotik.com/docs/spaces/ROS/pages/56459266/HotSpot+-+Captive+portal
- PPPoE: https://help.mikrotik.com/docs/spaces/ROS/pages/2031625/PPPoE
- MikroTik Wiki: https://wiki.mikrotik.com/

---

## RouterOS API vs CLI — Key Translation Rules

The RouterOS API closely follows CLI syntax with these transformations:

| CLI | API word format | Example |
|---|---|---|
| Command path with spaces | Replace spaces with `/` | `/ip address print` → `/ip/address/print` |
| Filter `?key=value` | Query word prefix `?` | `?disabled=false` |
| Filter exists | Query word no value | `?running` |
| Set a parameter | Attribute word `=key=value` | `=name=myqueue` |
| Limit returned fields | `.proplist` attribute | `=.proplist=name,address` |
| Reference by ID | `.id` attribute | `=.id=*1` |

**Reply word types:** `!re` = data row, `!done` = end of reply, `!trap` = device error, `!fatal` = connection error.

**Streaming vs one-shot commands:**
- One-shot (`Run*` in Go): `/print`, `/add`, `/set`, `/remove`, `/enable`, `/disable`
- Streaming (`Listen*` in Go): `/listen`, `print follow`, `print interval=N`, `/monitor-traffic`, `/ping`, `/torch`

---

## Common Print Flags (across all menus)

| Flag | Meaning |
|---|---|
| `X` | disabled |
| `D` | dynamic (auto-created) |
| `R` | running |
| `S` | slave |
| `A` | active |
| `I` | invalid / inactive |
| `*` | default item |
| `B` | blackhole / blocked |
| `C` | connected (route) |
| `U` | undoable / unreachable |

## Universal Print Subcommands

Every `print` supports these modifiers:

```
print brief              # Compact table
print detail             # All properties expanded
print terse              # Machine-readable compact
print value-list         # One property per line
print count-only         # Show count only
print where key=value    # Filter rows
print from=N             # Start from item N
print file=filename      # Save output to router file
print interval=Ns        # STREAMING — auto-refresh every N sec
print follow             # STREAMING — follow new entries (like tail -f)
print follow-only        # STREAMING — only show new entries
print stats              # RX/TX counters (interfaces, queues)
print stats-detail       # Detailed counters
print oid                # SNMP OID values
print as-value           # Key=value pairs
print bytes              # Byte counters (hotspot users)
print packets            # Packet counters
print status             # Status fields (hotspot active/host)
print without-paging     # No pause between pages
```

> **Streaming** modifiers (`interval`, `follow`, `follow-only`) produce continuous output. In Go, use `Listen*` for these — see go-routeros Skill.

---

## 1. QUEUE

### `/queue interface`
Maps interfaces to queue types. Read-only assignment.

```
/queue interface print
```
Output fields: `#`, `INTERFACE`, `QUEUE`, `ACTIVE-QUEUE`

Sample from device:
```
# INTERFACE   QUEUE                ACTIVE-QUEUE
0 ether1      only-hardware-queue  only-hardware-queue
5 l2tp-out1   no-queue             no-queue
```

Set queue type:
```
/queue interface set numbers=0 queue=default-small
```

---

### `/queue simple`
Per-IP bandwidth limiting. Packets match by src (upload) or dst (download) against `target`.

**Available operations:** `add`, `set`, `remove`, `enable`, `disable`, `move`, `reset-counters`, `reset-counters-all`, `comment`, `export`

**`add` parameters:**
```
name=             # Queue name (required)
target=           # IP/subnet or interface e.g. 192.168.1.10/32 or 192.168.1.0/24
dst=              # Destination filter (optional)
max-limit=        # Max speed rx/tx e.g. 10M/10M (MIR — Maximum Information Rate)
limit-at=         # Guaranteed speed rx/tx e.g. 1M/1M (CIR — Committed Information Rate)
burst-limit=      # Burst peak speed rx/tx e.g. 20M/20M
burst-threshold=  # Avg rate that triggers burst rx/tx e.g. 6M/6M
burst-time=       # Burst window rx/tx in seconds e.g. 8/8
bucket-size=      # Token bucket size (advanced burst)
total-max-limit=  # Bidirectional max (bits/s)
total-limit-at=   # Bidirectional guaranteed (bits/s)
total-burst-limit=
total-burst-threshold=
total-burst-time=
total-bucket-size=
total-queue=      # Queue kind for bidirectional e.g. pcq-upload-default
priority=         # 1–8 (1=highest)
queue=            # Queue kind rx/tx e.g. default-small/default-small
parent=           # Parent queue name for HTB hierarchy
packet-marks=     # Match mangle packet marks e.g. ppp-mark
place-before=     # Insert before item N
time=             # Time schedule e.g. 6h-22h,mon,tue,wed,thu,fri
comment=
disabled=yes/no
```

**`print` subcommands:** `brief`, `bytes`, `count-only`, `detail`, `file`, `follow`, `follow-only`, `from`, `interval`, `oid`, `packets`, `rate`, `stats`, `where`, `without-paging`

**`print stats`** — RX/TX bytes+packets+dropped per queue
**`print rate`** — current rate per queue

**Common examples:**
```
# Limit single client 10Mbps
/queue simple add name=client1 target=192.168.1.10/32 max-limit=10M/10M

# With burst: 50M max, 10M guaranteed, burst to 80M for 8s when below 40M
/queue simple add name=subnet1 target=192.168.230.0/24 \
    max-limit=50M/50M limit-at=10M/10M \
    burst-limit=80M/80M burst-threshold=40M/40M burst-time=8/8

# Reset counters
/queue simple reset-counters numbers=0
/queue simple reset-counters-all
```

> **From Mikrotik docs:** Simple queues process in sequence — each packet traverses every rule until a match. FastTrack bypasses simple queues; disable fasttrack if queues aren't applying.

---

### `/queue tree`
Advanced HTB queuing. Requires mangle packet marks. Used for global QoS.

**`add` parameters:**
```
name=
parent=           # global-in / global-out / global-total / interface-name / another-tree-queue
packet-mark=      # From /ip firewall mangle new-packet-mark=
max-limit=        # MIR (bits/s e.g. 10M)
limit-at=         # CIR (bits/s)
burst-limit=
burst-threshold=
burst-time=
bucket-size=
priority=         # 1–8
queue=            # Queue kind
comment=
disabled=yes/no
```

---

### `/queue type`
Available queue algorithms.

```
/queue type print
```
Built-in: `default`, `default-small`, `ethernet-default`, `wireless-default`, `hotspot-default`, `pcq-upload-default`, `pcq-download-default`, `only-hardware-queue`, `no-queue`

---

## 2. INTERFACE

### `/interface` (top level)

**`print`:**
```
/interface print
/interface print detail
/interface print stats
/interface print where type=ether
/interface print where running=yes
```
Output fields: `#`, `NAME`, `TYPE`, `ACTUAL-MTU`, `L2MTU`, `MAX-L2MTU`, `MAC-ADDRESS`
Flags: `D`=dynamic, `X`=disabled, `R`=running, `S`=slave

**`print detail`** additional: `last-link-up-time`, `last-link-down-time`, `link-downs`

Sample from device:
```
5  R  name="l2tp-out1" type="l2tp-out" mtu=1450 actual-mtu=1450
       last-link-down-time=mar/06/2026 10:31:12 last-link-up-time=mar/06/2026 10:31:28
       link-downs=1
```

**Monitor live traffic (STREAMING):**
```
/interface monitor-traffic ether1
/interface monitor-traffic ether1,ether2,ether3
```
Shows: `rx-bits-per-second`, `tx-bits-per-second`, `rx-packets-per-second`, `tx-packets-per-second`, `rx-drops`, `tx-drops`

---

### `/interface ethernet`

**`print`:**
```
/interface ethernet print
/interface ethernet print detail
/interface ethernet print stats
/interface ethernet print stats-detail
```
Output fields: `#`, `NAME`, `MTU`, `MAC-ADDRESS`, `ARP`, `SWITCH`

**`print detail`** fields: `name`, `default-name`, `mtu`, `l2mtu`, `mac-address`, `orig-mac-address`, `arp`, `arp-timeout`, `loop-protect`, `loop-protect-status`, `auto-negotiation`, `advertise`, `full-duplex`, `tx-flow-control`, `rx-flow-control`, `speed`, `bandwidth`, `switch`

Sample from device:
```
0 R name="ether1" mtu=1500 l2mtu=1520 mac-address=00:0C:42:73:88:36
     arp=enabled arp-timeout=auto loop-protect=default auto-negotiation=yes
     advertise=10M-half,10M-full,100M-half,100M-full,1000M-half,1000M-full
     full-duplex=yes tx-flow-control=off rx-flow-control=off speed=1Gbps
     bandwidth=unlimited/unlimited switch=switch1
```

**`set` parameters:**
```
name=
mtu=
mac-address=
arp=enabled/disabled/proxy-arp/reply-only
arp-timeout=auto/Ns
auto-negotiation=yes/no
advertise=10M-half,10M-full,100M-half,100M-full,1000M-half,1000M-full
full-duplex=yes/no
tx-flow-control=on/off/auto
rx-flow-control=on/off/auto
loop-protect=on/off/default
loop-protect-send-interval=  # e.g. 5s
loop-protect-disable-time=   # e.g. 5m
bandwidth=unlimited/N        # bps rx/tx
speed=10Mbps/100Mbps/1Gbps/auto
comment=
disabled=yes/no
```

Other operations: `blink` (LED), `cable-test` (detect faults), `monitor` (live stats), `reset-counters`, `reset-mac-address`

---

### `/interface bridge`

**`add` parameters:**
```
name=
mtu=
arp=enabled/disabled/proxy-arp/reply-only
protocol-mode=none/rstp/stp/mstp
priority=          # 0x0000–0xF000 (default 0x8000)
forward-delay=     # default 15s
hello-time=        # default 2s
max-message-age=   # default 20s
ageing-time=       # MAC expiry (default 5m)
igmp-snooping=yes/no
dhcp-snooping=yes/no
vlan-filtering=yes/no
comment=
disabled=yes/no
```

Sub-menus: `port`, `filter`, `nat`, `host` (MAC table), `vlan`, `msti`, `settings`, `mdb`

**`/interface bridge port add` parameters:**
```
interface=          # Physical port to bridge (required)
bridge=             # Bridge to add it to (required)
priority=           # 0x00–0xFF (default 0x80)
path-cost=          # STP cost (default 10)
pvid=               # Port VLAN ID (for vlan-filtering)
frame-types=        # admit-all / admit-only-untagged / admit-only-vlan-tagged
ingress-filtering=yes/no
edge=yes/no/auto    # PortFast
point-to-point=yes/no/auto
learn=yes/no/auto
hw=yes/no
horizon=
comment=
disabled=yes/no
```

---

### `/interface vlan`

**`add` parameters (from device):**
```
name=                    # e.g. vlan10 (required)
vlan-id=                 # 1–4094 (required)
interface=               # Parent interface (required)
mtu=
arp=enabled/disabled/proxy-arp/reply-only
arp-timeout=auto/Ns
loop-protect=on/off/default
loop-protect-send-interval=
loop-protect-disable-time=
use-service-tag=yes/no   # 802.1ad QinQ
comment=
disabled=yes/no
```

---

### `/interface pppoe-server`
Access Concentrator (ISP side). Sub-menu: `/interface pppoe-server server`

**`/interface pppoe-server server add` parameters:**
```
interface=           # Physical WAN interface (required)
service-name=        # Match client's service name
max-mtu=auto/N       # Usually auto → 1480
max-mru=auto/N
mrru=disabled/N      # Multilink PPP
authentication=      # mschap2,mschap1,chap,pap
default-profile=     # PPP profile for new sessions
keepalive-timeout=   # Seconds (0=never disconnect)
one-session-per-host=yes/no
max-sessions=        # 0=unlimited (Level 4 license = 200)
disabled=yes/no
```

---

### `/interface pppoe-client`

**`add` parameters:**
```
name=                # e.g. pppoe-wan (required)
interface=           # Physical WAN port (required)
user=                # Username from ISP
password=
service-name=        # Usually blank
ac-name=             # Usually blank
add-default-route=yes/no
default-route-distance=  # e.g. 1
use-peer-dns=yes/no
dial-on-demand=yes/no
profile=
max-mtu=auto/N
max-mru=auto/N
comment=
disabled=yes/no
```

---

### `/interface l2tp-client`

**`add` parameters:**
```
name=
connect-to=          # Remote server IP
user=
password=
profile=
use-ipsec=yes/no
ipsec-secret=
add-default-route=yes/no
default-route-distance=
keepalive-timeout=   # seconds
comment=
disabled=yes/no
```

---

## 3. PPP

### `/ppp secret`

**`print` (from device):**
```
/ppp secret print
/ppp secret print detail
/ppp secret print where service=pppoe
/ppp secret print where name~"user"
```
Fields: `#`, `NAME`, `SERVICE`, `CALLER-ID`, `PASSWORD`, `PROFILE`, `REMOTE-ADDRESS`
Flags: `X`=disabled

**`print detail`** additional: `routes`, `ipv6-routes`, `limit-bytes-in`, `limit-bytes-out`, `last-logged-out`

Sample from device:
```
0 name="ppp1" service=any caller-id="" password="" profile=default
  routes="" ipv6-routes="" limit-bytes-in=0 limit-bytes-out=0
  last-logged-out=jan/01/1970 00:00:00
```

**`add` parameters (from device):**
```
name=              # Username (required)
password=
service=           # any / pppoe / pptp / l2tp / ovpn / sstp (default: any)
profile=           # Profile from /ppp profile
local-address=     # Server IP or pool name
remote-address=    # Client IP or pool name
caller-id=         # Restrict by MAC (PPPoE) or calling number
routes=            # Static routes pushed to client
ipv6-routes=
limit-bytes-in=    # Download byte cap (0=unlimited)
limit-bytes-out=   # Upload byte cap (0=unlimited)
disabled=yes/no
comment=
```

> **From Mikrotik docs:** `/ppp secret` overrides `/ppp profile`. Exception: concrete IPs always win over pools when both define local/remote-address.

---

### `/ppp profile`

**`print` (from device):**
```
/ppp profile print
/ppp profile print detail
```
Flags: `*`=default

**`add` parameters (from device):**
```
name=                  # Profile name (required)
local-address=         # Server-side IP or pool
remote-address=        # Client IP or pool
dns-server=            # DNS IPs pushed to client (e.g. 8.8.8.8,8.8.4.4)
wins-server=
rate-limit=            # e.g. 10M/10M — creates queue for session
session-timeout=       # Auto-disconnect (e.g. 1h, 1d, 0=never)
idle-timeout=          # Idle disconnect (e.g. 5m, 0=never)
only-one=yes/no/default  # One session per username
incoming-filter=       # Firewall chain for incoming packets
outgoing-filter=       # Firewall chain for outgoing packets
parent-queue=
queue-type=
insert-queue-before=
change-tcp-mss=yes/no/default
use-compression=yes/no/default
use-encryption=yes/no/default
use-mpls=yes/no/default
use-upnp=yes/no/default
bridge=
bridge-learning=yes/no/default
bridge-horizon=
bridge-path-cost=
bridge-port-priority=
interface-list=
address-list=          # Add client IP to address-list on connect
on-up=                 # Script on connect
on-down=               # Script on disconnect
comment=
```

> **Note:** `incoming-filter`/`outgoing-filter` create dynamic jump rules into chain `ppp`. That chain must exist in `/ip firewall filter` first.

---

### `/ppp active`
Active PPP/PPPoE/L2TP sessions. **Read-only** (print/remove only).

**`print` (from device):**
```
/ppp active print
/ppp active print detail
/ppp active print stats              # rx/tx bytes+packets
/ppp active print follow interval=1  # STREAMING live monitor
/ppp active print where service=pppoe
```
Fields: `#`, `NAME`, `SERVICE`, `CALLER-ID`, `ADDRESS`, `UPTIME`, `ENCODING`
Flags: `R`=radius

Force disconnect:
```
/ppp active remove [find name="username"]
```

---

### `/ppp aaa`

**`print`:**
```
/ppp aaa print
```

**`set` parameters (from device):**
```
use-radius=yes/no
accounting=yes/no
interim-update=               # e.g. 5m
use-circuit-id-in-nas-port-id=yes/no
```

---

### `/ppp l2tp-secret`

**`print` (from device):**
```
/ppp l2tp-secret print
```
Fields: `#`, `ADDRESS`, `SECRET`

**`add` parameters (from device):**
```
address=   # Remote peer IP (required)
secret=    # Pre-shared key
comment=
```

---

## 4. IP DNS

### `/ip dns`

**`print` (from device):**
```
/ip dns print
```
Output:
```
servers:                 8.8.8.8,208.67.220.220,8.8.4.4,208.67.222.222
dynamic-servers:         192.168.47.1
use-doh-server:
verify-doh-cert:         no
allow-remote-requests:   yes
max-udp-packet-size:     4096
query-server-timeout:    2s
query-total-timeout:     10s
max-concurrent-queries:  100
max-concurrent-tcp-sessions: 20
cache-size:              2048KiB
cache-max-ttl:           1w
cache-used:              133KiB
```

**`set` parameters (from device):**
```
servers=                     # Upstream DNS comma-separated
allow-remote-requests=yes/no # Allow LAN clients to use this router as DNS
cache-size=                  # e.g. 2048KiB
cache-max-ttl=               # e.g. 1w, 1d
max-udp-packet-size=         # default 4096
query-server-timeout=        # e.g. 2s
query-total-timeout=         # e.g. 10s
max-concurrent-queries=
max-concurrent-tcp-sessions=
use-doh-server=              # DoH URL
verify-doh-cert=yes/no
```

Sub-menus: `cache` (print/flush), `static` (add/set/remove)

**`/ip dns static add`:**
```
name=        # Hostname e.g. myserver.local
address=     # IP to resolve to
type=A/CNAME/MX/TXT/NS/SRV/FWD
ttl=         # e.g. 1d
comment=
disabled=yes/no
```

---

## 5. IP POOL

### `/ip pool`

**`print` (from device):**
```
/ip pool print
```
Output:
```
# NAME          RANGES
0 dhcp_pool0    192.168.230.2-192.168.230.254
3 e2e-test-pool 10.88.0.1-10.88.0.10
```

**`add` parameters (from device):**
```
name=        # Pool name (required)
ranges=      # e.g. 192.168.1.10-192.168.1.254
             # Multiple: 192.168.1.10-192.168.1.100,192.168.1.150-192.168.1.200
next-pool=   # Overflow pool when exhausted
comment=
```

**`set` parameters (from device):** `name=`, `ranges=`, `next-pool=`, `comment=`

---

### `/ip pool used`
Allocated IPs from pools. **Read-only.**

**`print` (from device):**
```
/ip pool used print
```
Output:
```
POOL        ADDRESS           OWNER  INFO
dhcp_pool2  192.168.233.254   DHCP   88:88:88:88:87:88
```

---

## 6. IP ADDRESS

### `/ip address`

**`print` (from device):**
```
/ip address print
/ip address print detail
/ip address print where interface=ether1
```
Fields: `#`, `ADDRESS` (CIDR), `NETWORK`, `INTERFACE`
Flags: `X`=disabled, `I`=invalid, `D`=dynamic

Sample from device:
```
#   ADDRESS            NETWORK         INTERFACE
0   192.168.230.1/24   192.168.230.0   ether2
4 D 192.168.47.2/24    192.168.47.0    ether1      ← dynamic (DHCP client)
5 D 10.10.31.253/32    10.58.44.27     l2tp-out1   ← dynamic (L2TP)
```

**`add` parameters (from device):**
```
address=      # CIDR notation (required) e.g. 192.168.1.1/24
interface=    # Interface name (required)
network=      # Optional — auto-calculated
broadcast=    # Optional — auto-calculated
netmask=      # Alternative to CIDR prefix
comment=
disabled=yes/no
```

---

## 7. IP ROUTES

### `/ip route`

**`print` (from device):**
```
/ip route print
/ip route print detail
/ip route print terse
/ip route print where dst-address=0.0.0.0/0
```
Flags: `X`=disabled, `A`=active, `D`=dynamic, `C`=connect, `S`=static, `r`=rip, `b`=bgp, `o`=ospf, `B`=blackhole, `U`=unreachable, `P`=prohibit

Sample from device:
```
# FLAGS  DST-ADDRESS      PREF-SRC       GATEWAY        DISTANCE
0 ADS    0.0.0.0/0                       192.168.47.1   1
1 ADC    5.5.5.5/32       5.5.5.5        ether1         0
4  DC    192.168.230.0/24 192.168.230.1  ether2         255
```

**`add` parameters (from device):**
```
dst-address=   # Destination prefix (required) e.g. 0.0.0.0/0
gateway=       # Next-hop IP or interface name
distance=      # Admin distance 1–255 (default 1, lower = preferred)
check-gateway= # arp / ping / none — gateway reachability test
routing-mark=  # Policy routing mark from mangle
pref-src=      # Preferred source IP
type=          # unicast / blackhole / unreachable / prohibit
scope=         # Default 30
target-scope=  # Default 10
vrf-interface=
route-tag=
comment=
disabled=yes/no
```

**Verify route:**
```
/ip route check 8.8.8.8
```

---

### `/ip route nexthop`
Gateway reachability. **Read-only.**

**`print` (from device):**
```
/ip route nexthop print
```
Sample:
```
0 address=192.168.47.1 gw-state=reachable forwarding-nexthop="" interface="" scope=10 check-gateway=none
```

---

### `/ip route rule`
Policy-based routing.

**`print` (from device):**
```
/ip route rule print
```
Flags: `X`=disabled, `I`=inactive

Sample:
```
0 I action=lookup table=""
```

**`add` parameters (from device):**
```
src-address=    # Match source IP/subnet
dst-address=    # Match destination IP/subnet
interface=      # Match incoming interface
routing-mark=   # Match mangle routing mark
action=         # lookup / unreachable / prohibit / blackhole
table=          # Routing table name
place-before=
comment=
disabled=yes/no
```

---

## 8. IP FIREWALL

### `/ip firewall filter`

**`print` (from device):**
```
/ip firewall filter print
/ip firewall filter print detail
/ip firewall filter print chain=input
/ip firewall filter print stats
/ip firewall filter print where action=drop
```
Flags: `X`=disabled, `I`=invalid, `D`=dynamic

Sample from device:
```
0 X  ;;; place hotspot rules here
     chain=unused-hs-chain action=passthrough
1    ;;; BYPASS INTERNET POSITIF
     chain=input action=accept protocol=udp src-port=53
```

**`add` parameters:**
```
chain=           # input / forward / output or custom (required)
action=          # accept / drop / reject / log / passthrough / jump / return / tarpit
protocol=        # tcp / udp / icmp / gre / ipsec-esp / etc.
src-address=
dst-address=
src-port=
dst-port=
port=            # Match either src or dst port
in-interface=
out-interface=
in-interface-list=
out-interface-list=
src-address-list=
dst-address-list=
connection-state=  # new,established,related,invalid
connection-mark=
packet-mark=
layer7-protocol=
content=           # Match string in payload
limit=             # e.g. 10,5:packet
nth=               # Every Nth packet e.g. 2,1
random=            # 1–99 = percentage probability
psd=               # Port scan detection
hotspot=
icmp-options=      # Type:code e.g. 8:0
tcp-flags=
tcp-mss=
dscp=
ttl=
time=              # e.g. 6h-22h,mon-fri
fragment=yes/no
routing-mark=
ipsec-policy=
ipv4-options=
src-mac-address=
in-bridge-port=
out-bridge-port=
in-bridge-port-list=
out-bridge-port-list=
log=yes/no
log-prefix=
jump-target=       # Chain name when action=jump
reject-with=       # icmp-net-unreachable / icmp-host-unreachable / icmp-port-unreachable / tcp-reset
place-before=
comment=
disabled=yes/no
```

```
/ip firewall filter reset-counters numbers=0,1,2
/ip firewall filter reset-counters-all
```

---

### `/ip firewall nat`

**`print` (from device):**
```
/ip firewall nat print
/ip firewall nat print chain=srcnat
/ip firewall nat print stats
```
Flags: `X`=disabled, `I`=invalid, `D`=dynamic

Sample from device:
```
1    chain=srcnat action=masquerade log=no log-prefix=""
3    chain=srcnat action=masquerade out-interface=ether1 log=no log-prefix=""
4 X  ;;; DNS Premium United States
     chain=dstnat action=dst-nat to-addresses=208.67.220.220 to-ports=443
     protocol=udp dst-port=53 log=no log-prefix=""
```

**`add` parameters (from device) — common fields:**
```
chain=            # srcnat / dstnat (required)
action=           # masquerade / src-nat / dst-nat / netmap / redirect /
                  # passthrough / accept / drop / jump / return / log
src-address=
dst-address=
src-port=
dst-port=
port=
in-interface=
out-interface=
in-interface-list=
out-interface-list=
src-address-list=
dst-address-list=
protocol=
connection-state=
connection-mark=
packet-mark=
routing-mark=
layer7-protocol=
content=
limit=
nth=
random=
hotspot=
psd=
time=
fragment=yes/no
ipsec-policy=
src-mac-address=
to-addresses=     # Target IP/range for dst-nat/src-nat
to-ports=         # Target port for dst-nat/redirect
same-not-by-dst=yes/no
address-list=           # For add-src/dst-to-address-list actions
address-list-timeout=
log=yes/no
log-prefix=
jump-target=
place-before=
comment=
disabled=yes/no
```

Common patterns:
```
# Masquerade outbound (most common)
/ip firewall nat add chain=srcnat action=masquerade out-interface=ether1

# Port forward TCP 80 to internal server
/ip firewall nat add chain=dstnat protocol=tcp dst-port=80 \
    action=dst-nat to-addresses=192.168.1.10 to-ports=80

# Force DNS to router
/ip firewall nat add chain=dstnat protocol=udp dst-port=53 \
    action=redirect to-ports=53
```

---

### `/ip firewall mangle`

**`add` parameters (from device):**
```
chain=                 # prerouting / postrouting / input / output / forward (required)
action=                # mark-connection / mark-packet / mark-routing /
                       # change-ttl / change-mss / accept / drop / jump / return / log / passthrough
new-connection-mark=   # action=mark-connection
new-packet-mark=       # action=mark-packet
new-routing-mark=      # action=mark-routing
new-dscp=              # action=change-dscp
new-ttl=               # set/increment/decrement:N
new-mss=               # action=change-mss
new-priority=
passthrough=yes/no     # Continue to next rule after marking
route-dst=             # Force route target
routing-mark=
connection-mark=
packet-mark=
connection-state=      # new,established,related,invalid
connection-nat-state=  # srcnat,dstnat,no-nat
connection-bytes=
connection-rate=
per-connection-classifier=  # e.g. both-addresses:2/0 for load balancing
p2p=                        # P2P protocol detection
sniff-id=
sniff-target=
sniff-target-port=
# (all filter match fields also available: src-address, dst-address, protocol, ports, etc.)
comment=
disabled=yes/no
```

---

### `/ip firewall raw`
Pre-connection-tracking. Chains: `prerouting`, `output`. Actions include `notrack`.

```
# Bypass conntrack for high-throughput internal traffic
/ip firewall raw add chain=prerouting action=notrack \
    src-address=192.168.230.0/24 dst-address=192.168.230.0/24
```

Same match fields as filter. `add` parameters (from device): `action`, `address-list`, `address-list-timeout`, `comment`, `content`, `copy-from`, `disabled`, `dscp`, `dst-address`, `dst-address-list`, `dst-address-type`, `dst-limit`, `dst-port`, `fragment`, `hotspot`, `icmp-options`, `in-bridge-port`, `in-bridge-port-list`, `in-interface`, `in-interface-list`, `ingress-priority`, `ipsec-policy`, `ipv4-options`, `jump-target`, `limit`, `log`, `log-prefix`, `nth`, `out-bridge-port`, `out-bridge-port-list`, `out-interface`, `out-interface-list`, `packet-mark`, `packet-size`, `per-connection-classifier`, `place-before`, `port`, `priority`, `protocol`, `psd`, `random`, `src-address`, `src-address-list`, `src-address-type`, `src-mac-address`, `src-port`, `tcp-flags`, `tcp-mss`, `time`, `tls-host`, `ttl`, `chain`

---

### `/ip firewall address-list`

**`print` (from device):**
```
/ip firewall address-list print
```
Fields: `#`, `LIST`, `ADDRESS`, `CREATION-TIME`
Flags: `X`=disabled, `D`=dynamic

**`add` parameters (from device):** `address`, `comment`, `copy-from`, `disabled`, `timeout`, `list`

```
/ip firewall address-list add list=blocked-ips address=1.2.3.4
/ip firewall address-list add list=blocked-ips address=1.2.3.0/24 timeout=1h
```

---

### `/ip firewall connection`
Connection tracking table. **Read-only.**

```
/ip firewall connection print
/ip firewall connection print where src-address~"192.168.1"
/ip firewall connection print count-only
```

---

### `/ip firewall layer7-protocol`

**`add`:** `name=`, `regexp=`, `comment=`

---

### `/ip firewall service-port`
ALG helpers (ftp, tftp, sip, h323, pptp, etc.)

```
/ip firewall service-port print
/ip firewall service-port enable ftp
/ip firewall service-port disable sip
```

---

## 9. IP HOTSPOT

> **From Mikrotik docs:** HotSpot is a captive portal. Uses web-proxy internally. Requires IPv4. Enable in `/system device-mode` if needed.

**Setup wizard:** `/ip hotspot setup`

Sub-menus: `profile`, `user`, `user profile`, `active`, `host`, `cookie`, `ip-binding`, `walled-garden`, `service-port`

---

### `/ip hotspot profile`

**`print` (from device):**
```
/ip hotspot profile print
/ip hotspot profile print detail
```
Sample:
```
0 * name="default" hotspot-address=0.0.0.0 dns-name="" html-directory=hotspot
    html-directory-override="" rate-limit="" http-proxy=0.0.0.0:0
    smtp-server=0.0.0.0 login-by=cookie,http-chap http-cookie-lifetime=3d
    split-user-domain=no use-radius=no
```

**`add` parameters (from device):**
```
name=
hotspot-address=         # IP of hotspot server
dns-name=                # Captive portal hostname
html-directory=          # Servlet directory (default: hotspot)
html-directory-override=
http-proxy=              # Proxy IP:port
https-redirect=yes/no
smtp-server=             # Redirect SMTP
login-by=                # cookie,http-chap,http-pap,mac,radius,trial
mac-auth-mode=           # as-username / as-username-and-password
mac-auth-password=
http-cookie-lifetime=    # e.g. 3d
rate-limit=              # Global rate limit e.g. 5M/5M
ssl-certificate=
split-user-domain=yes/no
use-radius=yes/no
radius-accounting=yes/no
radius-interim-update=   # e.g. 5m
radius-default-domain=
radius-location-id=
radius-location-name=
radius-mac-format=
nas-port-type=
trial-uptime-limit=
trial-uptime-reset=
trial-user-profile=
```

---

### `/ip hotspot user`

**`print` (from device):**
```
/ip hotspot user print
/ip hotspot user print detail
/ip hotspot user print bytes
/ip hotspot user print packets
/ip hotspot user print where profile=1H
```
Fields: `#`, `SERVER`, `NAME`, `ADDRESS`, `PROFILE`, `UPTIME`
Flags: `*`=default, `X`=disabled, `D`=dynamic

Sample from device (mikhmon vouchers):
```
1  ;;; vc-mlysouc8-02.23.26
   TES452809   1H   0s
```

**`add` parameters (from device):**
```
name=               # Username / voucher code (required)
password=
profile=            # User profile from /ip hotspot user profile
server=             # Restrict to specific hotspot server
address=            # Lock to static IP
mac-address=        # Lock to MAC
limit-uptime=       # e.g. 1h, 1d, 30d
limit-bytes-in=     # Download cap (0=unlimited)
limit-bytes-out=    # Upload cap
limit-bytes-total=
routes=
email=
comment=            # Often used for expiry: YYYY/MM/DD:N (mikhmon format)
disabled=yes/no
```

---

### `/ip hotspot user profile`

**`add` parameters:**
```
name=                # Profile name (required)
rate-limit=          # e.g. 2M/2M or with burst: 1M/1M/2M/2M/8/8/1M/1M
shared-users=        # Max simultaneous sessions (default 1)
idle-timeout=        # e.g. 5m
keepalive-timeout=   # e.g. 2m
parent-queue=
queue-type=
insert-queue-before=
incoming-filter=
outgoing-filter=
address-list=        # Add to address-list on login
open-status-page=    # yes/no/http-login
status-autorefresh=
transparent-proxy=yes/no
on-login=            # Script on login
on-logout=           # Script on logout
```

---

### `/ip hotspot active`
**Read-only.** (print/remove)

**`print` (from device):**
```
/ip hotspot active print
/ip hotspot active print detail
/ip hotspot active print stats
/ip hotspot active print status
/ip hotspot active print follow interval=1   # STREAMING
/ip hotspot active print count-only
/ip hotspot active print where user~"TES"
```
Fields: `#`, `USER`, `ADDRESS`, `UPTIME`, `SESSION-TIME-LEFT`, `IDLE-TIMEOUT`
Flags: `R`=radius, `B`=blocked

Force logout: `/ip hotspot active remove [find user="username"]`

---

### `/ip hotspot host`
All hosts seen by gateway.

**`print` (from device):**
```
/ip hotspot host print
/ip hotspot host print detail
/ip hotspot host print status
/ip hotspot host print bytes
/ip hotspot host print packets
```
Fields: `#`, `MAC-ADDRESS`, `ADDRESS`, `TO-ADDRESS`, `SERVER`, `IDLE-TIMEOUT`
Flags: `S`=static, `H`=DHCP, `D`=dynamic, `A`=authorized, `P`=bypassed

`make-binding numbers=N` — Convert to static binding

---

### `/ip hotspot cookie`

**`print` (from device):**
```
/ip hotspot cookie print
```
Fields: `#`, `USER`, `DOMAIN`, `MAC-ADDRESS`, `EXPIRES-IN`
Flags: `M`=mac-cookie

Remove to force re-login: `/ip hotspot cookie remove [find user="username"]`

---

### `/ip hotspot ip-binding`

**`add` parameters:**
```
type=           # regular / bypassed / blocked
mac-address=
address=
server=
to-address=     # Assign static IP (type=regular)
comment=
disabled=yes/no
```

---

### `/ip hotspot walled-garden`

**`add` parameters:**
```
dst-host=       # e.g. *.google.com
dst-address=
dst-port=
method=         # GET/POST/all
server=
comment=
disabled=yes/no
```

---

## 10. SYSTEM LOGGING

### `/system logging`

**`print` (from device):**
```
/system logging print
```
Fields: `#`, `TOPICS`, `ACTION`, `PREFIX`
Flags: `X`=disabled, `I`=invalid, `*`=default

Sample from device:
```
0  * info       → memory
1  * error      → memory
3  * critical   → echo
4    account    → remote     prefix: api
5    hotspot,info,debug → disk  prefix: ->
```

**`add` parameters (from device):**
```
topics=      # info / error / warning / critical / debug / account / hotspot /
             # ppp / dhcp / firewall / l2tp / pptp / system / script /
             # wireless / ovpn / radius / dns / bridge / etc.
             # Combine: topics=hotspot,info,debug
             # Negate: !debug
action=      # Action name from /system logging action
prefix=      # Text prefix on each log line
disabled=yes/no
```

---

### `/system logging action`

**`print` (from device):**
```
/system logging action print
/system logging action print detail
```
Sample:
```
0 * name="memory"  target=memory  memory-lines=1000  memory-stop-on-full=no
1 * name="disk"    target=disk    disk-file-name="log"  disk-lines-per-file=1000  disk-file-count=2
2 * name="echo"    target=echo    remember=yes
3 * name="remote"  target=remote  remote=0.0.0.0  remote-port=514  src-address=0.0.0.0
                   bsd-syslog=no  syslog-time-format=bsd-syslog  syslog-facility=daemon  syslog-severity=auto
```

**`add` parameters (from device):**
```
name=                  # Action name (required)
target=                # memory / disk / echo / remote / email (required)

# target=memory:
memory-lines=          # Buffer size (default 1000)
memory-stop-on-full=yes/no

# target=disk:
disk-file-name=        # e.g. log
disk-lines-per-file=   # default 1000
disk-file-count=       # rotation count (default 2)
disk-stop-on-full=yes/no

# target=remote (syslog):
remote=                # Syslog server IP
remote-port=           # default 514
src-address=
bsd-syslog=yes/no
syslog-time-format=    # bsd-syslog / iso8601
syslog-facility=       # daemon / local0–local7 / kernel / user / etc.
syslog-severity=       # auto / debug / info / notice / warning / error / critical / alert / emergency

# target=email:
email-to=
email-start-tls=yes/no

remember=yes/no        # Keep messages after action (echo target)
comment=
```

---

## 11. SYSTEM SCHEDULER

**`print` (from device):**
```
/system scheduler print
/system scheduler print detail
```
Flags: `X`=disabled

Sample from device (mikhmon_expire_monitor):
```
name="mikhmon_expire_monitor"  start-date=feb/23/2026  start-time=00:00:00  interval=1m
on-event= :local montharray ("jan","feb","mar",...);
          :local date [/system clock get date];
          ...check hotspot user comment for expiry date YYYY/MM/DD:N format...
          if expired → remove user
```

**`add` parameters (from device):**
```
name=           # Scheduler name (required)
start-date=     # e.g. jan/01/2026 or startup
start-time=     # e.g. 00:00:00 / startup / never
interval=       # 0=run once, 30s, 1m, 1h, 1d, etc.
on-event=       # Inline script body OR script name from /system script
policy=         # ftp,reboot,read,write,policy,test,password,sniff,sensitive,romon
comment=
disabled=yes/no
```

Common patterns:
```
# Run every minute (expiry check like mikhmon)
/system scheduler add name=expire-check interval=1m \
    start-time=startup on-event=expire-script

# Daily midnight
/system scheduler add name=daily-backup start-time=00:00:00 \
    interval=1d on-event="/system backup save name=auto"

# Once on boot with delay
/system scheduler add name=boot-init start-time=startup \
    on-event=":delay 30; /system script run init-script"
```

---

## 12. SYSTEM SCRIPT

**`print` (from device):**
```
/system script print
/system script print detail
```
Fields: `name`, `owner`, `policy`, `dont-require-permissions`, `last-started`, `run-count`, `source`

Sample from device:
```
0 name="cache-update-trigger"  owner="admin"
  policy=ftp,reboot,read,write,policy,test,password,sniff,sensitive,romon
  dont-require-permissions=no  last-started=jan/20/2026 12:24:30  run-count=4
  source=
      :local resource "$resource"
      /tool fetch url="http://your-api:8080/api/cache/invalidate" \
          http-method=post http-data="..." as-value output=user
```

**`add` parameters (from device):**
```
name=
source=                     # RouterOS scripting language body
owner=                      # Username (default: current user)
policy=                     # Permission flags
dont-require-permissions=yes/no
comment=
```

Run: `/system script run my-script`

Sub-menus: `environment` (global vars), `job` (running scripts)

Policy flags: `ftp`, `reboot`, `read`, `write`, `policy`, `test`, `password`, `sniff`, `sensitive`, `romon`

---

## 13. SYSTEM (Misc)

### `/system resource`
**`print` (from device):**
```
/system resource print
/system resource print interval=1   # STREAMING
/system resource print oid
```
Output from device:
```
uptime:              2h21m24s
version:             6.49.11 (stable)
build-time:          Dec/08/2023 14:37:03
free-memory:         7.0MiB
total-memory:        32.0MiB
cpu:                 MIPS 24Kc V7.4
cpu-count:           1
cpu-frequency:       680MHz
cpu-load:            3%
free-hdd-space:      46.8MiB
total-hdd-space:     63.8MiB
write-sect-since-reboot: 564
write-sect-total:    20772182
bad-blocks:          0%
architecture-name:   mipsbe
board-name:          RB750G
platform:            G-Net
```

Sub-menus: `cpu`, `irq`, `pci`, `usb` (all print/monitor)

---

### `/system identity`
```
/system identity print   # → name: G-Net
/system identity set name=MyRouter
```

---

### `/system routerboard`
**`print` (from device):**
```
/system routerboard print
```
Output:
```
routerboard:      yes
model:            750G
serial-number:    228E01AED081
firmware-type:    ar7100
factory-firmware: 2.23
current-firmware: 2.23
upgrade-firmware: 6.49.11
```
Upgrade: `/system routerboard upgrade` then `/system reboot`

---

### `/system clock`
**`print` (from device):**
```
/system clock print
```
Output:
```
time:                23:38:26
date:                mar/06/2026
time-zone-autodetect: yes
time-zone-name:      Asia/Jakarta
gmt-offset:          +07:00
dst-active:          no
```

**`set` parameters (from device):** `time=`, `date=`, `time-zone-name=`, `time-zone-autodetect=yes/no`

`/system clock manual set`: `time-zone=`, `dst-delta=`, `dst-start=`, `dst-end=`

---

### `/system license`
**`print` (from device):**
```
/system license print
```
Output:
```
software-id: 75NH-QPYY
nlevel:      4
features:
```
Levels: 0=demo, 1, 3, 4(200 PPPoE), 5(500 PPPoE), 6(unlimited)

---

### `/system history`
```
/system history print
```
Flags: `U`=undoable, `R`=redoable, `F`=floating-undo

Undo: `:undo` at prompt

---

### `/system note`
```
/system note print
/system note set show-at-login=yes note="Authorized access only"
```

---

### Reboot / Shutdown
```
/system reboot
/system shutdown
```

---

## 14. LOG VIEWER

### `/log`

**`print` (from device):**
```
/log print
/log print detail          # time=, topics=, message=
/log print follow          # STREAMING — tail -f style
/log print follow-only     # STREAMING — only new entries
/log print where topics~"hotspot"
/log print where message~"login"
/log print terse
/log print count-only
```

Sample from device:
```
time=06:30:20 topics=system,error,critical
  message="router was rebooted without proper shutdown"

time=06:30:24 topics=l2tp,ppp,info
  message="l2tp-out1: initializing..."
```

Filter by severity: `/log info`, `/log warning`, `/log error`, `/log debug`

---

## 15. TOOLS

### `/tool profile`
```
/tool profile
/tool profile duration=10
```
Output from device:
```
NAME          CPU   USAGE
ethernet            0%
firewall            0%
winbox              1%
management          1.5%
hotspot             0%
queuing             0%
total               3%
```

---

### `ping` (root command)
STREAMING by default.
```
/ping 8.8.8.8
/ping 8.8.8.8 count=5 size=1400 interval=0.5
/ping 8.8.8.8 src-address=192.168.1.1
/ping 192.168.1.1 arp-ping=yes interface=ether2
```
Parameters: `count=`, `interface=`, `routing-table=`, `src-address=`, `size=`, `ttl=`, `dscp=`, `interval=`, `arp-ping=yes/no`

Sample from device:
```
SEQ HOST     SIZE TTL TIME   STATUS
  0 8.8.8.8   56  116 24ms
sent=6 received=6 packet-loss=0% min-rtt=22ms avg-rtt=23ms max-rtt=24ms
```

---

### `/tool traceroute`
```
/tool traceroute 8.8.8.8
/tool traceroute 8.8.8.8 src-address=192.168.1.1 count=3
```

---

### `/tool bandwidth-test`
```
/tool bandwidth-test 192.168.1.1 user=admin password=admin \
    direction=both duration=10 protocol=tcp
```

---

### `/tool torch`
STREAMING realtime per-flow traffic (like iftop).
```
/tool torch interface=ether1
/tool torch interface=ether1 src-address=192.168.1.0/24
/tool torch interface=ether1 port=80 protocol=tcp
```

---

### `/tool ping-speed`
```
/tool ping-speed address=8.8.8.8
```
Output from device: `current: 0bps`, `average: 0bps`

---

### `/tool netwatch`
```
/tool netwatch add host=8.8.8.8 interval=10s timeout=1s \
    up-script=":log info \"gateway up\"" \
    down-script=":log warning \"gateway down\""
```

---

### `/tool fetch`
```
# Download file
/tool fetch url="http://example.com/file.txt" dst-path=flash/file.txt

# Webhook POST (as used in device scripts)
/tool fetch url="http://your-api:8080/api/cache/invalidate" \
    http-method=post \
    http-data="{\"resource\":\"user\",\"action\":\"add\"}" \
    as-value output=user
```
Parameters: `url=`, `dst-path=`, `http-method=get/post/put`, `http-data=`, `http-header-field=`, `mode=http/https/ftp/tftp`, `as-value`, `output=none/file/user`, `user=`, `password=`

---

### `/tool mac-server` / `/tool mac-telnet`
```
/tool mac-server set allowed-interface-list=all
/tool mac-telnet 00:0C:42:73:88:36
```

---

### `/tool graphing`
```
/tool graphing interface add interface=ether1 store-every=5min allow-address=0.0.0.0/0
/tool graphing resource add allow-address=0.0.0.0/0
```

---

### `/tool e-mail`
```
/tool e-mail set server=smtp.gmail.com port=587 tls=starttls \
    from=router@example.com user=user password=pass

/tool e-mail send to=admin@example.com subject="Alert" body="Router rebooted"
```

---

### `/tool ip-scan` / `/tool mac-scan`
```
/tool ip-scan interface=ether2
/tool mac-scan interface=ether2 duration=5
```

---

## 16. IP DHCP SERVER

### `/ip dhcp-server`

**`add` parameters:**
```
name=
interface=          # Interface to serve (required)
address-pool=       # Pool from /ip pool
lease-time=         # e.g. 1d, 12h
authoritative=      # yes/after-2sec/no
use-radius=yes/no
conflict-detection=yes/no
comment=
disabled=yes/no
```

---

### `/ip dhcp-server lease`
```
/ip dhcp-server lease print
/ip dhcp-server lease print detail
/ip dhcp-server lease make-static numbers=0
```
Fields: `ADDRESS`, `MAC-ADDRESS`, `CLIENT-ID`, `HOST-NAME`, `STATUS`, `EXPIRES-AFTER`

---

### `/ip dhcp-server network`

**`add` parameters:**
```
address=          # Subnet CIDR (required) e.g. 192.168.1.0/24
gateway=          # Default gateway IP
dns-server=       # Comma-separated DNS IPs
netmask=
domain=
wins-server=
ntp-server=
next-server=      # PXE/TFTP boot server
boot-file-name=
```

---

### `/ip dhcp-client`
```
/ip dhcp-client print
/ip dhcp-client add interface=ether1 disabled=no \
    add-default-route=yes use-peer-dns=yes
/ip dhcp-client release numbers=0
/ip dhcp-client renew numbers=0
```

---

## CLI → API Path Mapping

> For Go code to send these commands, see the **go-routeros** Skill.

| CLI menu | API command prefix |
|---|---|
| `/queue simple print` | `/queue/simple/print` |
| `/queue tree print` | `/queue/tree/print` |
| `/interface print` | `/interface/print` |
| `/interface ethernet print` | `/interface/ethernet/print` |
| `/interface bridge print` | `/interface/bridge/print` |
| `/interface vlan print` | `/interface/vlan/print` |
| `/ppp secret print` | `/ppp/secret/print` |
| `/ppp profile print` | `/ppp/profile/print` |
| `/ppp active print` | `/ppp/active/print` |
| `/ip pool print` | `/ip/pool/print` |
| `/ip pool used print` | `/ip/pool/used/print` |
| `/ip address print` | `/ip/address/print` |
| `/ip route print` | `/ip/route/print` |
| `/ip route nexthop print` | `/ip/route/nexthop/print` |
| `/ip firewall filter print` | `/ip/firewall/filter/print` |
| `/ip firewall nat print` | `/ip/firewall/nat/print` |
| `/ip firewall mangle print` | `/ip/firewall/mangle/print` |
| `/ip firewall address-list print` | `/ip/firewall/address-list/print` |
| `/ip firewall connection print` | `/ip/firewall/connection/print` |
| `/ip hotspot user print` | `/ip/hotspot/user/print` |
| `/ip hotspot active print` | `/ip/hotspot/active/print` |
| `/ip hotspot host print` | `/ip/hotspot/host/print` |
| `/ip hotspot cookie print` | `/ip/hotspot/cookie/print` |
| `/ip hotspot profile print` | `/ip/hotspot/profile/print` |
| `/ip dns print` | `/ip/dns/print` |
| `/ip dhcp-server lease print` | `/ip/dhcp-server/lease/print` |
| `/system resource print` | `/system/resource/print` |
| `/system identity print` | `/system/identity/print` |
| `/system logging print` | `/system/logging/print` |
| `/system scheduler print` | `/system/scheduler/print` |
| `/system script print` | `/system/script/print` |
| `/system clock print` | `/system/clock/print` |
| `/log print` | `/log/print` |

---

## Streaming Commands — Use `Listen*` in Go

> All commands below must use `Listen*` (not `Run*`) when called via go-routeros. See the **go-routeros** Skill for implementation.

| CLI command | Go `Listen*` arguments |
|---|---|
| `/ping address=8.8.8.8` | `/ping`, `=address=8.8.8.8` |
| `/tool torch interface=ether1` | `/tool/torch`, `=interface=ether1` |
| `/log print follow` | `/log/print`, `=follow=` |
| `/log print follow-only` | `/log/print`, `=follow-only=` |
| `/interface monitor-traffic ether1` | `/interface/monitor-traffic`, `=interface=ether1` |
| `print interval=N` (any menu) | add `=interval=N` |
| `print follow` (any menu) | add `=follow=` |
| `/ppp active print follow` | `/ppp/active/print`, `=follow=` |
| `/ip hotspot active print follow` | `/ip/hotspot/active/print`, `=follow=` |
| `/system resource print interval=1` | `/system/resource/print`, `=interval=1` |
| `/ip firewall address-list listen` | `/ip/firewall/address-list/listen` |
| `/queue simple print follow` | `/queue/simple/print`, `=follow=` |

> **Always** set `c.Queue = 100` before `Listen*` calls to buffer incoming sentences during data bursts.