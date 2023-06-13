<script setup lang="ts">
import MemosView from "@/views/HomeView/MemosView.vue";
import { useRouter, useRoute } from "vue-router";
import { ref } from "vue";

const isMove = ref(false);

console.log("Page Route: ",useRoute().query)
let page = useRoute().query['page'];
if (page == "edit") {
    isMove.value = true;
    console.log(useRoute().query['fileKey'])
    useRouter().push({
        path: "edit",
        query: {
            "fileKey": useRoute().query['fileKey']
        }
    })
}

</script>

<template>
  <main>
    <MemosView v-if="!isMove"/>
  </main>
</template>
