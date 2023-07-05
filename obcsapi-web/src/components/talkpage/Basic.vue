<script setup lang="ts">
import marked from "marked";

defineProps<{
    mdText: string, // dayKey
}>();

function markdown(text: string) :string {
  return marked(text || '', {
        breaks: true
  })
}
</script>

<template>
    <div class="message1" v-html="markdown(mdText)"></div>
</template>
<style scoped>
.chat-messages {
  overflow-y: auto;
  margin-bottom: 10px;
  overflow: hidden;
  /* 添加此样式以隐藏溢出内容 */
}

.chat-messages :deep() img {
  max-width: 100%;
  display: block;
}

.message1 {
  background-color: rgba(255, 255, 255, 0.05);
  padding: 10px;
  border-radius: 5px;
  margin-bottom: 5px;
}

.message2 {
  background-color: rgba(255, 255, 255, 0.01);
  padding: 10px;
  border-radius: 5px;
  margin-bottom: 5px;
}


.chat-input {
  display: flex;
  gap: 10px;
}
</style>
  