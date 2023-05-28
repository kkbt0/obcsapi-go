<script setup lang="ts">
import { ref } from "vue";
import { type UploadFileInfo, NUpload } from "naive-ui";

const url: string = localStorage.getItem("host") + "/api/v1/upload";
const token: string = localStorage.getItem("token") || "";
const previewFileList = ref<UploadFileInfo[]>([])
const emit = defineEmits(['upload-callback'])


const handleUploadFinish = ({ file, event }: {
  file: UploadFileInfo
  event?: ProgressEvent
}) => {
  // console.log(file)
  // console.log(file.url)
  // console.log(event)
  file.url = JSON.parse((event?.target as XMLHttpRequest).response).data.url
  if (isImageFile(file.url || "unknown")) {
    emit('upload-callback', `![${file.name}](${file.url})`)
  } else {
    emit('upload-callback', `[${file.name}](${file.url})`)
  }
  // inputText.value += `\n![${file.name}](${file.url})\n`;
  return file
}

function isImageFile(filename: string): boolean {
  // 将文件名后缀转换为小写字母
  var extension = filename.toLowerCase().split('.').pop() || "unknown";
  // 图片文件的常见后缀
  var imageExtensions = ['jpg', 'jpeg', 'gif', 'png', 'svg', 'ico', 'webp'];
  // 检查后缀是否在图片后缀列表中
  return imageExtensions.includes(extension);
}



</script>
<template>
  <n-upload :action="url" :default-file-list="previewFileList" :headers="{ 'Authorization': token }"
    @finish="handleUploadFinish" list-type="image-card">
    点击上传
  </n-upload>
</template>