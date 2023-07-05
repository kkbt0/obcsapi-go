<script setup lang="ts">
import { ref, watch } from 'vue';
import marked from "marked";

const emit = defineEmits(["web-command"]);

const props = defineProps<{
    inobj: any, // dayKey
}>();


const obj = ref(props.inobj);

watch(() => props.inobj, (newVal) => {
    obj.value = newVal;
});


function ClickHandler(command: string) {
    emit("web-command", obj.value ,command);
}

function markdown(text: string) :string {
  return marked(text || '', {
        breaks: true
  })
}

</script>
<template>
    <div class="message1" v-html="markdown(obj.data.parts[0].text)"></div>
    <n-space justify="space-around">
        <n-button quaternary type="info" @click="ClickHandler(obj.data.parts[1].command)">{{ obj.data.parts[1].text
        }}</n-button>
        <n-button quaternary type="success" @click="ClickHandler(obj.data.parts[2].command)">{{ obj.data.parts[2].text }}</n-button>
    </n-space>
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
  