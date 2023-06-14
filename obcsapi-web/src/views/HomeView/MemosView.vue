<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRouter, useRoute } from "vue-router";
import { NSpace } from "naive-ui"
import MemosInput from "@/components/obcsapi/MemosInput.vue";
import MemosList from "@/components/obcsapi/MemosList.vue";
import { ObcsapiTestJwt } from "@/api/obcsapi";

const router = useRouter();

onMounted(() => {
    console.log(`Initiating`)
    ObcsapiTestJwt().then(text => {
        if (text != "hello") {
            router.push("/login");
        }
    }).catch(err => {
        console.log(err)
    })
})


</script>

<template>
    <br>
    <n-space vertical>
        <MemosInput />
        <MemosList />
    </n-space>
</template>

<style scoped></style>
