<script setup lang="ts">
import { ref } from "vue";
import { NInput, NSpace, NButton } from "naive-ui";
import { ObcsapiPostMemos } from "@/api/obcsapi";
import { memosData } from "@/stores/memos";
import MemosUpload  from "@/components/obcsapi/MemosUpload.vue";

const inputText = ref("");
const memos = memosData();
const showUpload = ref(false);


function sendMemos() {
  if (inputText.value) {
    ObcsapiPostMemos("", 9999, "", inputText.value).then(data => {
      memos.setMap(data.date, data)
      inputText.value = "";
      showUpload.value = false;
      window.$message.success("Suceess Send");
    }).catch(e => {
      console.log(e);
    });
  } else {
    window.$message.warning("Empty message")
  }
}

function imgUrlDeal(text:string) {
  inputText.value += `\n${text}\n`;
}

</script>

<template>
  <n-space vertical>
    <n-input v-model:value="inputText" type="textarea" class="memos-input" placeholder="Memos"
      :autosize="{ minRows: 3 }" />
    <n-space justify="space-between">
      <n-button quaternary type="info" @click="showUpload=!showUpload">📌</n-button>
      <n-button quaternary type="primary" @click="sendMemos">🚀</n-button>
    </n-space>
    <MemosUpload v-if="showUpload" @upload-callback="imgUrlDeal" />
  </n-space>
</template>

<style scoped>
.memos-input {
  min-width: 300px;
}
</style>
