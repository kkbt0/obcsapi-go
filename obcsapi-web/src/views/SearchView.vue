<script setup lang="ts">
import { NList, NListItem, NThing, NEmpty, NMention,NBackTop } from 'naive-ui';
import { ref, onMounted } from 'vue';
import { ObcsapiSerchKvCache } from '@/api/obcsapi';
import { useRoute, useRouter } from 'vue-router';
import { LocalSetting } from '@/stores/setting';

const inputText = ref("");
const result = ref();
const router = useRouter();

onMounted(() => {
    if (useRoute().query['key'] != undefined) {
        console.log("Search: ",useRoute().query['key']);
        inputText.value = <string>useRoute().query['key'];
        searchServer();
    }
})

function searchServer() {
    inputText.value = inputText.value.trim();
    ObcsapiSerchKvCache(inputText.value).then(res => {
        if (res.status != 200) {
            window.$message.warning(`${res.status}`)
        } else {
            res.json().then(obj => {
                result.value = obj;
            });
        }
    })
}
function clickSearch(fileKey: string) {
    console.log(fileKey);
    router.push({
        path: "edit",
        query: {
            "fileKey": fileKey
        }
    })
}

function gotoEdit() {
    router.push({
        path: "edit",
    })
}
</script>
<template>
    <br>
    <n-space vertical>
        <n-mention v-model:value="inputText" :prefix="['#']" :options="LocalSetting().mention" />
        <n-space justify="end">
            <n-button @click="searchServer">Search</n-button>
        </n-space>
        <n-back-top :right="25" />
        <div v-if="!result || result.length == 0">
            <n-empty description="空">
                <template #extra>
                    <n-button size="small" @click="gotoEdit">
                        返回
                    </n-button>
                </template>
            </n-empty>
        </div>
        <div v-else>
            <a>Result: {{ result.length }} files</a>
        </div>
        <n-list bordered>
            <n-list-item v-for="(val, index) in result" :key="index" class="search-item">
                <n-thing :title="val.filekey" @click="clickSearch(val.filekey)">
                    <div v-html="val.content.replaceAll(inputText, `<a>${inputText}</a>`)"></div>
                </n-thing>
            </n-list-item>
        </n-list>
    </n-space>
</template>
<style scoped>
.search-item :deep() div {
    word-break: break-all;
}
</style>