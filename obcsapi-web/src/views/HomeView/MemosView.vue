<script setup lang="ts">
import { NSpace } from "naive-ui"
import MemosInput from "@/components/obcsapi/MemosInput.vue";
import MemosList from "@/components/obcsapi/MemosList.vue";



</script>

<template>
    <br>
    <n-space vertical>
        <MemosInput />
        <MemosList />
    </n-space>
</template>

<style scoped></style>
