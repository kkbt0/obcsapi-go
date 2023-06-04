<script setup lang="ts">
import { reactive, nextTick, ref } from "vue";
import { NList, NListItem, NScrollbar, NSpace } from "naive-ui"
import { memosData } from "@/stores/memos";
import MemosEdit from "@/components/obcsapi/MemosEdit.vue"

const memosIndexList = reactive(memosData().memosMap);

// 自动加载更多
let scroll: any;
//窗口可视范围高度
let clientHeight: any;
//文档内容实际高度  总高度（包括超出视窗的溢出部分）
let scrollHeight: number;
const isLoading = ref(false);

upDateScroll();
async function upDateScroll() {
    await nextTick();
    scrollInit();
}

function scrollInit() {
    scroll = document.getElementById("scrollId");
    clientHeight = scroll.innerHeight || Math.min(scroll.clientHeight, scroll.clientHeight);
    let children = scroll.querySelectorAll(".child");
    scrollHeight = 0;
    children.forEach(function (child: any) {
        scrollHeight += child.offsetHeight;
    });
}

function LoadMoreMemosList() {
    isLoading.value = true;
    memosData().waitMoreMemos(20).then(() => {
        scrollInit()
        isLoading.value = false;
    }).catch((e) => {
        console.log(e);
        isLoading.value = false;
    });
}


function scrollEvent(e: any) {
    //滚动条滚动距离
    let scrollTop = e.target.scrollTop;
    if (clientHeight + scrollTop >= scrollHeight) {
        if (isLoading.value) {
            return;
        }
        console.log("scroll end");
        LoadMoreMemosList();
    }
}
</script>

<template>
    <n-scrollbar style="max-height: 80vh" id="scrollId" @scroll="scrollEvent">
        <div v-if="memosData().memosIndexList" v-for="(dayKey, key1) in memosData().memosIndexList" :key="key1"
            class="child">
            <!-- 每天的列表 -->
            <n-list bordered v-if="memosIndexList.get(dayKey) != undefined">
                <!-- 列表中的 Memos -->
                <div v-for="(memosShowText, key2) in memosIndexList.get(dayKey).md_show_text.slice().reverse()" :key="key2">
                    <n-list-item v-if="memosShowText.trim() != ''" class="n-list-item-custom">
                        <MemosEdit :memosShowText="memosShowText"
                            :memosRaw="memosIndexList.get(dayKey).md_text[memosIndexList.get(dayKey).md_text.length - key2 - 1]"
                            :dayKey="dayKey" :line="memosIndexList.get(dayKey).md_text.length - key2 - 1" />
                    </n-list-item>
                    <div />
                </div>
            </n-list>
        </div>
        <n-space justify="space-around">
        <n-button quaternary type="primary" @click="LoadMoreMemosList" :disabled="isLoading">{{ isLoading ? 'Loading More' :
            'Load More' }}</n-button>
    </n-space>
    </n-scrollbar>

</template>

<style scoped>
.n-list :deep() .n-list-item-custom {
    padding: 10px 12px;
}
</style>