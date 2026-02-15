import { ref, onUnmounted } from 'vue'

/**
 * 语音识别 composable
 * 优先使用浏览器 Web Speech API，无需后端
 */
export function useSpeechRecognition() {
  const isListening = ref(false)
  const transcript = ref('')
  const isSupported = ref(false)
  const error = ref<string | null>(null)

  let recognition: SpeechRecognition | null = null

  const SpeechRecognitionAPI =
    window.SpeechRecognition || (window as any).webkitSpeechRecognition

  if (SpeechRecognitionAPI) {
    isSupported.value = true
    recognition = new SpeechRecognitionAPI()
    recognition.continuous = false
    recognition.interimResults = true
    recognition.lang = 'zh-CN'

    recognition.onresult = (event: SpeechRecognitionEvent) => {
      let text = ''
      for (let i = event.resultIndex; i < event.results.length; i++) {
        text += event.results[i][0].transcript
      }
      transcript.value = text
    }

    recognition.onend = () => {
      isListening.value = false
    }

    recognition.onerror = (event: any) => {
      error.value = event.error
      isListening.value = false
    }
  }

  function start(lang = 'zh-CN') {
    if (!recognition) {
      error.value = 'Speech recognition not supported'
      return
    }
    transcript.value = ''
    error.value = null
    recognition.lang = lang
    recognition.start()
    isListening.value = true
  }

  function stop() {
    recognition?.stop()
    isListening.value = false
  }

  onUnmounted(() => {
    stop()
  })

  return { isListening, transcript, isSupported, error, start, stop }
}
