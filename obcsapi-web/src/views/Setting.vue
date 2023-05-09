<script setup lang="ts">
import MemosConfig from "@/components/obcsapi/MemosConfig.vue";
import { useRouter } from "vue-router";
import {  ObcsapiTestMail } from "@/api/obcsapi";

const router = useRouter()
function clearCache() {
    localStorage.removeItem("mainMdList")
    window.$message.info("Clearing")
}
function reLogin() {
    localStorage.removeItem("token")
    window.$message.info("注销")
    router.push("/")
}
function sendMail() {
    ObcsapiTestMail().then( res => {
        window.$message.info(res)
    }).catch( e => {
        window.$message.info(e)
    })
}
</script>
<template>
    <MemosConfig />
    <n-space >
        <n-button @click="sendMail" quaternary>测试邮件</n-button>
        <n-button @click="clearCache" type="info" quaternary>清除缓存</n-button>
        <n-button @click="reLogin" type="warning" quaternary>注销</n-button>
    </n-space>
</template>