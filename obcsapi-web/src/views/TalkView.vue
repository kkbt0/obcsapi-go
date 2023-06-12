<script setup lang="ts">
import { ref, type Ref } from 'vue';
import { NButton, NScrollbar, NMention } from 'naive-ui';
import { ObcsapiTalk } from "@/api/obcsapi";
import { TalkStore } from "@/stores/talk";
import { LocalSetting } from "@/stores/setting"

import marked from "marked";
import MemosUpload from "@/components/obcsapi/MemosUpload.vue";

const messages = TalkStore().messages;
const newMessage = ref('');
const scrollbarRef: Ref<any> = ref(null);
const contentRef: Ref<any> = ref(null);
const talkInputRef = ref();
const showUpload = ref(false);

function sendMessage() {
  if (newMessage.value == '') {
    window.$message.warning("Empty");
    return;
  }
  messages.push("I: " + newMessage.value);
  ObcsapiTalk(newMessage.value).then(text => {
    showUpload.value = false;
    messages.push("O: " + text);
    scrollToBottom()
  }).catch(e => {
    console.log(e);
    window.$message.warning("Err Save: " + e);
  })
  newMessage.value = '';
}

const scrollToBottom = () => {
  // console.log(scrollbarRef.value);
  const contentHeight = contentRef.value.clientHeight;
  const scrollY = contentHeight - window.innerHeight * 0.75;;
  // console.log(contentHeight, scrollY);
  scrollbarRef.value.scrollTo(0, scrollY);
};

function addTags() {
  newMessage.value += '#'
  talkInputRef.value?.focus();
}

function imgUrlDeal(text: string) {
  newMessage.value += `\n${text}\n`;
}

function markdown(text: string) :string {
  return marked(text || '', {
        breaks: true
    })
}

</script>

<template>
  <div>
    <div>Talk</div>
    <div class="chat-input">
      <n-mention type="textarea" :autosize="{ minRows: 2 }" :options="LocalSetting().mention" v-model:value="newMessage"
        placeholder=":~$" :prefix="['#']" ref="talkInputRef"/>
    </div>
    <n-space justify="space-between">
      <n-button quaternary type="info" @click="showUpload = !showUpload">📌</n-button>
      <n-button quaternary type="info" @click="addTags">🏷️</n-button>
      <n-button quaternary type="info" @click="sendMessage">🚀</n-button>
    </n-space>
    <MemosUpload v-if="showUpload" @upload-callback="imgUrlDeal" />
    <n-scrollbar style="max-height: 75vh" ref="scrollbarRef">
      <div class="chat-messages" ref="contentRef">
        <div v-for="(message, index) in messages.slice().reverse()" :key="index">
          <div v-if="message.substring(0, 2) == 'I:'" class="message1" v-html="markdown(message)"></div>
          <div v-else class="message2" v-html="markdown(message)"></div>
        </div>
      </div>
    </n-scrollbar>
  </div>
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
  