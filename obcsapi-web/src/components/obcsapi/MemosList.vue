<script setup lang="ts">
import { NList, NListItem } from "naive-ui"
import { memosData } from "@/stores/memos";
import MemosEdit from "@/components/obcsapi/MemosEdit.vue"


// 去除空白的数组
function arrFilter(strList: Array<string>): Array<string> {
    return strList.filter((item) => {
        return item.trim() != '';
    });
}


</script>

<template>
    <div v-if="memosData().memosIndexList" v-for="(dayKey, key1) in memosData().memosIndexList" :key="key1">
        <!-- 每天的列表 -->
        <n-list bordered v-if="memosData().memosMap.get(dayKey) != undefined">
            <!-- 列表中的 Memos -->
            <n-list-item v-for="(memosShowText, key2) in memosData().memosMap.get(dayKey).md_show_text.slice().reverse()"
                :key="key2">
                <MemosEdit 
                    :memosShowText="memosShowText"
                    :memosRaw="memosData().memosMap.get(dayKey).md_text[memosData().memosMap.get(dayKey).md_text.length - key2 - 1]"
                    :dayKey="dayKey"
                    :line="memosData().memosMap.get(dayKey).md_text.length - key2 - 1" 
                />
            </n-list-item>
        </n-list>
    </div>
</template>

<style scoped></style>
