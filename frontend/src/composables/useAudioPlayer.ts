import { ref, onUnmounted } from 'vue'

/**
 * 音频播放器 composable
 * 支持流式音频播放 (收到 PCM/MP3 chunk 立即播放)
 */
export function useAudioPlayer() {
  const isPlaying = ref(false)
  const enabled = ref(true)
  let audioContext: AudioContext | null = null
  let nextStartTime = 0
  const queue: AudioBuffer[] = []
  let playing = false

  function getContext(): AudioContext {
    if (!audioContext) {
      audioContext = new AudioContext()
    }
    return audioContext
  }

  /** 播放 MP3/OGG 二进制数据 */
  async function playChunk(data: ArrayBuffer) {
    if (!enabled.value) return
    const ctx = getContext()
    const buffer = await ctx.decodeAudioData(data.slice(0))
    queue.push(buffer)
    if (!playing) {
      scheduleNext()
    }
  }

  function scheduleNext() {
    const ctx = getContext()
    const buffer = queue.shift()
    if (!buffer) {
      playing = false
      isPlaying.value = false
      return
    }
    playing = true
    isPlaying.value = true
    const source = ctx.createBufferSource()
    source.buffer = buffer
    source.connect(ctx.destination)
    const startTime = Math.max(ctx.currentTime, nextStartTime)
    source.start(startTime)
    nextStartTime = startTime + buffer.duration
    source.onended = () => {
      scheduleNext()
    }
  }

  /** 停止播放并清空队列 */
  function stop() {
    queue.length = 0
    if (audioContext) {
      audioContext.close()
      audioContext = null
    }
    nextStartTime = 0
    playing = false
    isPlaying.value = false
  }

  /** 重置时间 (用于新轮次对话) */
  function reset() {
    queue.length = 0
    nextStartTime = 0
    playing = false
    isPlaying.value = false
  }

  onUnmounted(() => {
    stop()
  })

  return { isPlaying, enabled, playChunk, stop, reset }
}
