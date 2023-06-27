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
    ObcsapiTestJwt().then(res => {
        if (res.code !=200) {
            router.push("/login");
            window.$message.error(`${res.code}`);
        }
    }).catch(err => {
        window.$message.error(err);
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
