<script setup lang="ts">
import { reactive } from "vue";
import { NList, NListItem } from "naive-ui"
import { memosData } from "@/stores/memos";
import MemosEdit from "@/components/obcsapi/MemosEdit.vue"

const memosIndexList = reactive(memosData().memosMap)

</script>

<template>
    <div v-if="memosData().memosIndexList" v-for="(dayKey, key1) in memosData().memosIndexList" :key="key1">
        <!-- 每天的列表 -->
        <n-list bordered v-if="memosIndexList.get(dayKey) != undefined">
            <!-- 列表中的 Memos -->
            <n-list-item v-for="(memosShowText, key2) in memosIndexList.get(dayKey).md_show_text.slice().reverse()"
                :key="key2">
                <MemosEdit 
                    :memosShowText="memosShowText"
                    :memosRaw="memosIndexList.get(dayKey).md_text[memosIndexList.get(dayKey).md_text.length - key2 - 1]"
                    :dayKey="dayKey"
                    :line="memosIndexList.get(dayKey).md_text.length - key2 - 1" 
                />
            </n-list-item>
        </n-list>
    </div>
</template>

<style scoped></style>
