<script setup lang="ts">
import { onMounted } from "vue";
import { useRouter } from "vue-router";
import { NSpace } from "naive-ui"
import MemosInput from "@/components/obcsapi/MemosInput.vue";
import MemosList from "@/components/obcsapi/MemosList.vue";

import { ObcsapiTestJwt } from "@/api/obcsapi";

const router = useRouter();

onMounted(() => {
    console.log(`Initiating`)
    ObcsapiTestJwt().then(text => {
        if (text!="hello") {
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
    <main>
        <n-space vertical>
            <MemosInput />
            <MemosList />
        </n-space>
    </main>
</template>

<style scoped></style>
