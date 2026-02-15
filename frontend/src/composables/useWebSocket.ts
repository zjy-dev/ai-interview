import { ref, onUnmounted } from 'vue'

export interface WSMessage {
  type: 'text' | 'audio' | 'status' | 'error' | 'evaluation'
  data: any
}

/**
 * WebSocket composable - 面试实时交互
 */
export function useWebSocket() {
  const connected = ref(false)
  const error = ref<string | null>(null)
  let ws: WebSocket | null = null
  const handlers = new Map<string, ((data: any) => void)[]>()

  function connect(url: string) {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const wsUrl = `${protocol}//${window.location.host}${url}`
    const token = localStorage.getItem('token')

    ws = new WebSocket(`${wsUrl}?token=${token}`)
    ws.binaryType = 'arraybuffer'

    ws.onopen = () => {
      connected.value = true
      error.value = null
    }

    ws.onmessage = (event) => {
      if (event.data instanceof ArrayBuffer) {
        emit('audio', event.data)
        return
      }
      try {
        const msg: WSMessage = JSON.parse(event.data)
        emit(msg.type, msg.data)
      } catch {
        emit('text', event.data)
      }
    }

    ws.onerror = () => {
      error.value = 'WebSocket error'
    }

    ws.onclose = () => {
      connected.value = false
    }
  }

  function send(data: string | object) {
    if (!ws || ws.readyState !== WebSocket.OPEN) return
    ws.send(typeof data === 'string' ? data : JSON.stringify(data))
  }

  function on(event: string, handler: (data: any) => void) {
    if (!handlers.has(event)) handlers.set(event, [])
    handlers.get(event)!.push(handler)
  }

  function emit(event: string, data: any) {
    handlers.get(event)?.forEach((h) => h(data))
  }

  function disconnect() {
    ws?.close()
    ws = null
    connected.value = false
    handlers.clear()
  }

  onUnmounted(() => {
    disconnect()
  })

  return { connected, error, connect, send, on, disconnect }
}
