<script setup lang="ts">
import { reactive, ref, onMounted } from "vue";
import { NList, NListItem, NScrollbar, NSpace } from "naive-ui"
import { memosData } from "@/stores/memos";
import MemosEdit from "@/components/obcsapi/MemosEdit.vue"
import { ObcsapiGetMemos } from "@/api/obcsapi";

const memosIndexList = reactive(memosData().memosMap)
const memos = memosData();
const loadMemosCount = ref(0); // 取负数 然后 waitMoreMemos()

onMounted(() => {
    waitMoreMemos(20)
})

async function moreMemos() {
    await ObcsapiGetMemos(memos.dayBefore).then(data => {
        console.log(data.date)
        memos.addDaily(data.date, data)
        loadMemosCount.value += data.md_text.length // 统计 memos 数量
        memos.dayBefore -= 1;
    }).catch(err => {
        console.log(err)
    })
}

async function waitMoreMemos(needNum: number) {
    console.log(`Loading more ${needNum} memos`)
    let maxRequest = 5
    let tem = loadMemosCount.value + needNum
    while (loadMemosCount.value < tem && maxRequest > 0) {
        await moreMemos()
        maxRequest -= 1
    }

}

function LoadMoreMemosList() {
    waitMoreMemos(20)
}

</script>

<template>
    <n-scrollbar style="max-height: 75vh">
        <div v-if="memosData().memosIndexList" v-for="(dayKey, key1) in memosData().memosIndexList" :key="key1">
            <!-- 每天的列表 -->
            <n-list bordered v-if="memosIndexList.get(dayKey) != undefined">
                <!-- 列表中的 Memos -->
                <div v-for="(memosShowText, key2) in memosIndexList.get(dayKey).md_show_text.slice().reverse()" :key="key2">
                    <n-list-item v-if="memosShowText.trim() != ''">
                            <MemosEdit :memosShowText="memosShowText"
                            :memosRaw="memosIndexList.get(dayKey).md_text[memosIndexList.get(dayKey).md_text.length - key2 - 1]"
                            :dayKey="dayKey" :line="memosIndexList.get(dayKey).md_text.length - key2 - 1" />
                    </n-list-item>
                    <div />
                </div>
            </n-list>
        </div>
    </n-scrollbar>
    <n-space justify="end">
        <n-button quaternary type="primary" @click="LoadMoreMemosList">Load More</n-button>
    </n-space>
</template>

<style scoped></style>
