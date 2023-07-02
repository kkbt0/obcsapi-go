<script setup lang="ts">
import MemosView from "@/views/HomeView/MemosView.vue";
import { useRouter, useRoute } from "vue-router";
import { ref } from "vue";
import { LoginByCode, SayHello } from "@/router/jump";

const isMove = ref(false);
const router = useRouter();

console.log("Page Route: ", useRoute().query)

// check code
LoginByCode();
// 监测登录状态
SayHello()


// 跳转 edit 页
let page = useRoute().query['page'];
if (page == "edit") {
  isMove.value = true;
  console.log(useRoute().query['fileKey'])
  router.push({
    path: "edit",
    query: {
      "fileKey": useRoute().query['fileKey']
    }
  })
}

</script>

<template>
  <main>
    <MemosView v-if="!isMove" />
  </main>
</template>
