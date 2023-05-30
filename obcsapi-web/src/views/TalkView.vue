<script setup lang="ts">
import { ref, type Ref, onMounted } from 'vue';
import { NButton, NInput, NScrollbar,NMention } from 'naive-ui';
import { ObcsapiTalk } from "@/api/obcsapi";
import { TalkStore } from "@/stores/talk";
import { LocalSetting } from "@/stores/setting"
import { range } from 'lodash';

const messages = TalkStore().messages;
const newMessage = ref('');
const scrollbarRef: Ref<any> = ref(null);
const contentRef: Ref<any> = ref(null);
const mentionList: Ref<Array<any>> = ref([]); // 提示词

onMounted(() => {
  LocalSetting().mention.forEach(val => {
    console.log(val);
    mentionList.value.push({ label: val, value: val });
  })
})

function sendMessage() {
  if (newMessage.value == '') {
    window.$message.warning("Empty");
    return;
  }
  messages.push("I: " + newMessage.value);
  ObcsapiTalk(newMessage.value).then(text => {
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

</script>

<template>
  <div>
    <div>Talk</div>
    <div class="chat-input">
      <n-mention type="textarea" :autosize="{ minRows: 2 }" :options="mentionList" v-model:value="newMessage" placeholder=":~$"  :prefix="['#']"/>
      <n-button @click="sendMessage">输入</n-button>
    </div>
    <n-scrollbar style="max-height: 75vh" ref="scrollbarRef">
      <div class="chat-messages" ref="contentRef">
        <div v-for="(message, index) in messages.slice().reverse()" :key="index">
          <div v-if="message.substring(0, 2) == 'I:'" class="message1" v-html="message"></div>
          <div v-else class="message2" v-html="message"></div>
        </div>
      </div>
    </n-scrollbar>
  </div>
</template>
  

  
<style>
.chat-messages {
  overflow-y: auto;
  margin-bottom: 10px;
  overflow: hidden;
  /* 添加此样式以隐藏溢出内容 */
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
  