import type { MessageApiInjection } from "naive-ui/lib/message/src/MessageProvider"
declare global {
  interface Window {
      $message: MessageApiInjection
  }
}
