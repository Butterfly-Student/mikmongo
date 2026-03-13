---
name: go-routeros
description: Write Go code that connects to and controls Mikrotik RouterOS devices using the go-routeros/routeros/v3 library. Use when tasks involve RouterOS API, Mikrotik automation, network device management, or Go code that reads/writes router config.
dependencies: go>=1.21, github.com/go-routeros/routeros/v3
---

# RouterOS Go Client — `go-routeros/routeros/v3`

## Overview

This Skill provides complete guidance for writing Go programs that communicate with Mikrotik devices using the RouterOS API. Apply it whenever the task involves:
- Connecting to a Mikrotik router from Go code
- Reading router configuration or statistics
- Setting/modifying router config programmatically
- Subscribing to realtime events from a router
- Building monitoring tools, automation scripts, or network management apps in Go

**Import path:** `github.com/go-routeros/routeros/v3`
**Install:** `go get github.com/go-routeros/routeros/v3`

---

## Core Concepts

### Two Operation Modes

**Sync mode (default):** One command at a time, blocking. Best for simple scripts and CLI tools.

**Async mode:** Concurrent commands via multiplexing. Enable with `c.Async()`. Best for realtime listeners and concurrent goroutines.

### Two Command Patterns

**Run** — Send a command and wait for the complete reply. Use for queries and mutations (print, add, remove, set).

**Listen** — Subscribe to a streaming command that pushes updates (e.g. `/listen`, `/monitor`). Returns a channel of sentences. Non-blocking.

---

## Connecting to a Device

```go
import "github.com/go-routeros/routeros/v3"

// Plain TCP (port 8728)
c, err := routeros.Dial("192.168.1.1:8728", "admin", "password")

// With timeout
c, err := routeros.DialTimeout("192.168.1.1:8728", "admin", "password", 10*time.Second)

// With context (cancellable)
c, err := routeros.DialContext(ctx, "192.168.1.1:8728", "admin", "password")

// TLS (port 8729)
c, err := routeros.DialTLS("192.168.1.1:8729", "admin", "password", nil)

// TLS with context
c, err := routeros.DialTLSContext(ctx, "192.168.1.1:8729", "admin", "password", nil)

defer c.Close()
```

> All `Dial*` functions automatically call `Login()`. Use `NewClient(rwc)` only when you need to provide a custom `io.ReadWriteCloser` (e.g. for testing/mocking).

---

## Running Commands (Blocking)

```go
// Variadic form — most common
reply, err := c.Run("/ip/address/print")

// Slice form — useful when building commands dynamically
reply, err := c.RunArgs([]string{"/ip/address/print", "?interface=ether1"})

// With context (recommended for production)
reply, err := c.RunContext(ctx, "/system/resource/print")
reply, err := c.RunArgsContext(ctx, []string{"/interface/print", "?running=true"})
```

### Reading Reply Data

```go
// reply.Re — slice of *proto.Sentence (the !re rows)
// reply.Done — the final !done sentence
for _, sentence := range reply.Re {
    fmt.Println(sentence.Map["name"], sentence.Map["address"])
}
```

---

## RouterOS Command Syntax

RouterOS API uses a different format than the CLI:

| Purpose | Format | Example |
|---|---|---|
| Command | `/path/command` | `/ip/address/print` |
| Filter (boolean) | `?key=value` | `?disabled=false` |
| Filter (exists) | `?key` | `?running` |
| Set parameter | `=key=value` | `=interface=ether1` |
| Limit fields returned | `=.proplist=k1,k2` | `=.proplist=name,address` |
| Reference item by ID | `=.id=*1` | Use `.id` from a previous print |

**Example — add firewall rule:**
```go
_, err := c.Run(
    "/ip/firewall/filter/add",
    "=chain=input",
    "=protocol=tcp",
    "=dst-port=22",
    "=action=accept",
)
```

**Example — remove an item using its ID:**
```go
// Step 1: get ID
reply, _ := c.Run("/ip/address/print", "?interface=ether1", "=.proplist=.id")
id := reply.Re[0].Map[".id"]

// Step 2: remove by ID
_, err := c.Run("/ip/address/remove", "=.id="+id)
```

---

## Listening to Realtime Streams

Use `Listen*` for commands that stream updates over time.

```go
// Set buffer size before listening (prevents blocking on burst data)
c.Queue = 100

l, err := c.ListenArgsContext(ctx, []string{"/ip/firewall/address-list/listen"})
if err != nil {
    log.Fatal(err)
}

// Cancel after timeout in background
go func() {
    time.Sleep(30 * time.Second)
    l.Cancel()
}()

// Read updates
for sentence := range l.Chan() {
    fmt.Println(sentence.Map)
}

// ALWAYS check error after channel closes
if err := l.Err(); err != nil {
    log.Fatal(err)
}
```

> **Critical:** Always call `l.Err()` after `l.Chan()` is closed to catch any errors that occurred during streaming.

---

## Async Mode (Concurrent Commands)

Enable when running multiple goroutines that each send commands:

```go
errCh := c.Async() // or c.AsyncContext(ctx)

// Handle async errors in a separate goroutine
go func() {
    if err := <-errCh; err != nil {
        log.Fatal(err)
    }
}()

// Now safe to call Run/Listen from multiple goroutines
var wg sync.WaitGroup
for _, iface := range interfaces {
    wg.Add(1)
    go func(name string) {
        defer wg.Done()
        r, _ := c.Run("/interface/print", "?name="+name)
        fmt.Println(r.Re[0].Map)
    }(iface)
}
wg.Wait()
```

---

## Error Handling

```go
import "errors"

reply, err := c.Run("/some/command")
if err != nil {
    var deviceErr *routeros.DeviceError
    var unknownErr *routeros.UnknownReplyError

    switch {
    case errors.As(err, &deviceErr):
        // Error from RouterOS itself (!trap or !fatal)
        // e.g. permission denied, unknown command
        fmt.Println("RouterOS error:", deviceErr.Sentence.Map["message"])
    case errors.As(err, &unknownErr):
        // Unexpected reply word — library or protocol issue
        fmt.Println("Unknown reply:", unknownErr.Sentence)
    default:
        // Network/connection error
        log.Fatal(err)
    }
}
```

---

## Logging / Debug

```go
// Uses slog.Handler — pass any compatible handler
handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
})
c.SetLogHandler(handler)
// Will now log all API sentences sent and received
```

---

## Complete Working Example — Print Interface Stats

```go
package main

import (
    "fmt"
    "log"
    "strings"

    "github.com/go-routeros/routeros/v3"
)

func main() {
    c, err := routeros.Dial("192.168.1.1:8728", "admin", "password")
    if err != nil {
        log.Fatal(err)
    }
    defer c.Close()

    props := "name,rx-byte,tx-byte,rx-packet,tx-packet"
    reply, err := c.Run(
        "/interface/print",
        "?disabled=false",
        "?running=true",
        "=.proplist="+props,
    )
    if err != nil {
        log.Fatal(err)
    }

    for _, re := range reply.Re {
        for _, p := range strings.Split(props, ",") {
            fmt.Printf("%s=%s\t", p, re.Map[p])
        }
        fmt.Println()
    }
}
```

---

## Complete Working Example — Realtime Monitor with Context

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "time"

    "github.com/go-routeros/routeros/v3"
)

func main() {
    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
    defer cancel()

    c, err := routeros.DialContext(ctx, "192.168.1.1:8728", "admin", "password")
    if err != nil {
        log.Fatal(err)
    }
    defer c.Close()

    c.Queue = 100

    l, err := c.ListenArgsContext(ctx, []string{"/ip/firewall/address-list/listen"})
    if err != nil {
        log.Fatal(err)
    }

    // Auto-cancel after 60s
    go func() {
        time.Sleep(60 * time.Second)
        l.Cancel()
    }()

    for {
        select {
        case <-ctx.Done():
            return
        case sen, ok := <-l.Chan():
            if !ok {
                if err := l.Err(); err != nil {
                    log.Fatal(err)
                }
                return
            }
            fmt.Printf("Update: %v\n", sen.Map)
        }
    }
}
```

---

## Quick Reference: Method Cheatsheet

| Task | Method |
|---|---|
| Connect (plain) | `Dial`, `DialTimeout`, `DialContext` |
| Connect (TLS) | `DialTLS`, `DialTLSTimeout`, `DialTLSContext` |
| Run command (blocking) | `Run`, `RunArgs`, `RunContext`, `RunArgsContext` |
| Stream updates | `Listen`, `ListenArgs`, `ListenArgsContext`, `ListenArgsQueueContext` |
| Enable concurrency | `Async()`, `AsyncContext(ctx)` |
| Stop a stream | `l.Cancel()`, `l.CancelContext(ctx)` |
| Read stream data | `l.Chan()` (returns `<-chan *proto.Sentence`) |
| Check stream error | `l.Err()` |
| Enable debug logging | `c.SetLogHandler(slog.Handler)` |
| Close connection | `c.Close()` |

---

## When to Use Each Pattern

**Use `Run`** when you need a one-shot result: print config, add a rule, remove an address, set a value. Always blocking — call from main goroutine or with context for timeout.

**Use `Listen`** when the router pushes updates continuously: monitoring traffic, watching firewall events, tracking DHCP leases. Non-blocking — read from a goroutine.

**Use `Async`** when you need to send multiple commands concurrently from different goroutines. Required before calling `Listen` inside a goroutine that also calls `Run`.

**Always set `c.Queue`** before `Listen` to avoid dropped updates during traffic bursts. A value of 100 is a safe default.