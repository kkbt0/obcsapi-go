<script setup lang="ts">
import { ref } from "vue";
import { NInput, NSpace, NButton, type UploadFileInfo, NUpload } from "naive-ui";
import { ObcsapiPostMemos } from "@/api/obcsapi";
import { memosData } from "@/stores/memos";

const inputText = ref("");
const memos = memosData();
const showUpload = ref(false);
const token = localStorage.getItem("token") || "";
const previewFileList = ref<UploadFileInfo[]>([])


function sendMemos() {
  if (inputText.value) {
    ObcsapiPostMemos("", 9999, "", inputText.value).then(data => {
      memos.memosMap.set(data.date, data)
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


function showUploadChange() {
  showUpload.value = !showUpload.value;
}

const handleUploadFinish = ({ file, event }: {
  file: UploadFileInfo
  event?: ProgressEvent
}) => {
  // console.log(event)
  file.url = JSON.parse((event?.target as XMLHttpRequest).response).data.url
  // console.log(file)
  // console.log(file.url)
  inputText.value += `\n![${file.name}](${file.url})\n`;
  return file
}

</script>

<template>
  <n-space vertical>
    <n-input v-model:value="inputText" type="textarea" class="memos-input" placeholder="Input Memos"
      :autosize="{ minRows: 3 }" />
    <n-space justify="space-between">
      <n-button ghost type="info" @click="showUploadChange">Img</n-button>
      <n-button ghost type="primary" @click="sendMemos">Send</n-button>
    </n-space>
    <n-upload v-if="showUpload" action="http://localhost:8900/api/v1/upload" :default-file-list="previewFileList"
      :headers="{ 'Authorization': token }" @finish="handleUploadFinish" list-type="image-card">
      点击上传
    </n-upload>
  </n-space>
</template>

<style scoped>
.memos-input {
  min-width: 300px;
}
</style>
