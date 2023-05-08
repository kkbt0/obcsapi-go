<script setup lang="ts">
import { ref } from "vue";
import { NInput, NSpace, NButton } from "naive-ui";
import { ObcsapiPostMemos } from "@/api/obcsapi";
import { memosData } from "@/stores/memos";

const inputText = ref("");
const memos = memosData();

function sendMemos() {
  if (inputText.value) {
    console.log(inputText.value);
    ObcsapiPostMemos("", 9999, "", inputText.value).then(data => {
      memos.memosMap.set(data.date, data)
      inputText.value = "";
      window.$message.success("Suceess Send");
    }).catch(e => {
      console.log(e);
    });
  } else {
    window.$message.warning("Empty message")
  }

}
</script>

<template>
  <n-space vertical>
    <n-input v-model:value="inputText" type="textarea" class="memos-input" placeholder="Input Memos"
      :autosize="{ minRows: 3 }" />
    <n-space>
      <n-button @click="sendMemos">Send</n-button>
    </n-space>
  </n-space>
</template>

<style scoped>
.memos-input {
  min-width: 300px;
}
</style>
