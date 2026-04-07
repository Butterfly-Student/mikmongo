import { useEffect, useCallback, useRef, useState } from 'react'

type PingStatus = 'connecting' | 'connected' | 'error' | 'unavailable'

interface PingState {
    latency: string
    status: PingStatus
}

const MAX_RECONNECT_ATTEMPTS = 5
const RECONNECT_DELAY_MS = 3000

export function usePing(routerId: string | null): PingState {
    const [pingState, setPingState] = useState<PingState>({
        latency: '-- ms',
        status: 'unavailable',
    })

    const wsRef = useRef<WebSocket | null>(null)
    const reconnectTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null)
    const reconnectAttemptsRef = useRef(0)
    const isMountedRef = useRef(true)
    const currentRouterIdRef = useRef<string | null>(null)

    const cleanup = useCallback(() => {
        if (reconnectTimeoutRef.current) {
            clearTimeout(reconnectTimeoutRef.current)
            reconnectTimeoutRef.current = null
        }
        if (wsRef.current) {
            wsRef.current.onclose = null
            wsRef.current.onerror = null
            wsRef.current.onmessage = null
            wsRef.current.onopen = null
            wsRef.current.close()
            wsRef.current = null
        }
    }, [])

    const connect = useCallback((id: string) => {
        cleanup()
        reconnectAttemptsRef.current = 0

        const attemptConnect = () => {
            if (!isMountedRef.current || currentRouterIdRef.current !== id) return

            const internalKey = import.meta.env.VITE_INTERNAL_KEY as string
            // Always use window.location.origin so WS goes through the Vite proxy in dev
            // and through the reverse proxy in production — never connect directly to the backend port
            const wsBase = window.location.origin.replace(/^http/, 'ws')
            const wsUrl = `${wsBase}/api/v1/routers/${id}/monitor/ws/ping?address=8.8.8.8&token=${encodeURIComponent(internalKey)}`

            if (isMountedRef.current) {
                setPingState({ latency: '...', status: 'connecting' })
            }

            const ws = new WebSocket(wsUrl)
            wsRef.current = ws

            ws.onopen = () => {
                if (!isMountedRef.current || currentRouterIdRef.current !== id) {
                    ws.close()
                    return
                }
                reconnectAttemptsRef.current = 0
                setPingState({ latency: '-- ms', status: 'connected' })
            }

            ws.onmessage = (event) => {
                if (!isMountedRef.current || currentRouterIdRef.current !== id) return
                try {
                    const data = JSON.parse(event.data)
                    // Server sent an error — stop retrying (e.g. router not found)
                    if (data.error) {
                        setPingState({ latency: '-- ms', status: 'error' })
                        reconnectAttemptsRef.current = MAX_RECONNECT_ATTEMPTS
                        return
                    }
                    const latencyMs = data.timeMs ?? data.time ?? data.latency ?? data.avg_rtt ?? null
                    if (latencyMs !== null) {
                        setPingState({
                            latency: `${Math.round(Number(latencyMs))} ms`,
                            status: 'connected',
                        })
                    }
                } catch {
                    // ignore malformed messages
                }
            }

            ws.onerror = () => {
                if (!isMountedRef.current || currentRouterIdRef.current !== id) return
                setPingState({ latency: '-- ms', status: 'error' })
            }

            ws.onclose = () => {
                if (!isMountedRef.current || currentRouterIdRef.current !== id) return
                setPingState({ latency: '-- ms', status: 'error' })

                if (reconnectAttemptsRef.current < MAX_RECONNECT_ATTEMPTS) {
                    reconnectAttemptsRef.current += 1
                    reconnectTimeoutRef.current = setTimeout(attemptConnect, RECONNECT_DELAY_MS)
                }
            }
        }

        attemptConnect()
    }, [cleanup])

    useEffect(() => {
        isMountedRef.current = true

        if (!routerId) {
            cleanup()
            currentRouterIdRef.current = null
            setPingState({ latency: '-- ms', status: 'unavailable' })
            return
        }

        currentRouterIdRef.current = routerId
        connect(routerId)

        return () => {
            cleanup()
        }
    }, [routerId, connect, cleanup])

    useEffect(() => {
        return () => {
            isMountedRef.current = false
        }
    }, [])

    return pingState
}
