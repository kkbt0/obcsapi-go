<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRouter, useRoute } from "vue-router";
import { NSpace } from "naive-ui"
import MemosInput from "@/components/obcsapi/MemosInput.vue";
import MemosList from "@/components/obcsapi/MemosList.vue";
import { ObcsapiTestJwt } from "@/api/obcsapi";

const router = useRouter();
const isMove = ref(false);

console.log(useRoute().query)
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

onMounted(() => {
    console.log(`Initiating`)
    ObcsapiTestJwt().then(text => {
        if (text != "hello") {
            router.push("/login");
        }
    }).catch(err => {
        router.push("/login");
        console.log(err)
    })
})


</script>

<template>
    <br>
    <main v-if="!isMove">
        <n-space vertical>
            <MemosInput />
            <MemosList />
        </n-space>
    </main>
</template>

<style scoped></style>
