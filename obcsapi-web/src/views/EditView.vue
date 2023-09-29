<script lang="ts" setup >
import { ref, type Ref } from 'vue'
import { type TreeOption, NTree, NSpace } from 'naive-ui'
import { ObcsapiListFile } from "@/api/obcsapi"
import Editor from '@/components/obcsapi/Editor.vue';
import DailyCal from '@/components/obcsapi/DailyCal.vue';
import { LocalSetting } from '@/stores/setting';
import { useRoute , useRouter } from "vue-router";

const showMode = ref(0);
const fileKey = ref("edit.md");
const router = useRouter();
fileKey.value = <string>useRoute().query['fileKey'] || LocalSetting().recentEditList[0];

const dataTree: Ref<TreeOption[]> = ref([
    getRecentList()
])
// 用于点击事件判断的文件列表
const fileKeyList: Ref<string[]> = ref(LocalSetting().recentEditList.slice());
// 从 Local 中初始化文件树
ReBuildTree(LocalSetting().allFileKeyList);

function ReGetObcsapiListFile() {
    ObcsapiListFile().then(list => {
        ReBuildTree(list);
        LocalSetting().allFileKeyList = list;
        localStorage.setItem('AllFileKeyList', JSON.stringify(LocalSetting().allFileKeyList));
    })
}

function ReBuildTree(list: any) {
    // ['xx.md']
    list.forEach((i: string) => {
        fileKeyList.value.push(i);
    })
    dataTree.value = [ getRecentList() ]; // 初始化文件树
    buildTree(list).forEach(item => {
        dataTree.value.push(item);
    })
}

const updatePrefixWithExpaned = (
    _keys: Array<string | number>,
    _option: Array<TreeOption | null>,
    meta: {
        node: TreeOption | null
        action: 'expand' | 'collapse' | 'filter'
    }
) => {
    if (!meta.node) return
    switch (meta.action) {
        case 'expand':
            break
        case 'collapse':
            break
    }
}
const nodeProps = ({ option }: { option: TreeOption }) => {
    return {
        onClick() {
            if (!option.children && !option.disabled) {
                window.$message.info('Edit: ' + option.label)
                fileKeyList.value.forEach(val => {
                    // 因为 Ob 不允许同名文件存在，所以可以这样
                    if (val.endsWith(option.key?.toString()!)) {
                        fileKey.value = val;
                    }
                })
            }
        }
    }
}

function buildTree(data: string[]): TreeOption[] {
    const tree: TreeOption = {
        key: 'File',
        label: 'File',
        children: []
    };
    data.filter(item => item.endsWith('.md')).forEach(item => {
        const parts = item.split('/');
        let node = tree;
        for (let i = 0; i < parts.length - 1; i++) {
            const part = parts[i];
            let child = node.children?.find(c => c.key === part);
            if (!child) {
                child = {
                    key: part,
                    label: part,
                    children: []
                };
                node.children?.push(child);
            }
            node = child;
        }
        node.children?.push({
            key: parts[parts.length - 1],
            label: parts[parts.length - 1]
        });
    });
    return tree.children!;
}

function getRecentList(): TreeOption {
    const tree: TreeOption = {
        key: 'RecentEdit',
        label: 'RecentEdit',
        children: []
    };
    LocalSetting().recentEditList.forEach(item => {
        tree.children?.push({
            key: item,
            label: item,
        })
    })
    return tree;
}

function CalClicks(infileKey: string) {
    console.log("Load:",infileKey+".md");
    fileKey.value = infileKey+".md";
    showMode.value = 0;
}

function goSearchPage() {
    router.push("/search");
}

</script>
<template>
    <n-space vertical>
        <DailyCal v-show="showMode==2" @cal-click="CalClicks"/>
        <n-tree v-show="showMode==1" block-line expand-on-click :data="dataTree" :node-props="nodeProps"
            :on-update:expanded-keys="updatePrefixWithExpaned" />
        <a v-show="showMode==1" @click="ReGetObcsapiListFile">ReGetObcsapiListFile</a>
        <Editor :fileKey="fileKey" style="height: 85vh;" />
        <n-space justify="space-between">
            <a @click="showMode =2">Cal</a>
            <a @click="showMode =1">FileList</a>
            <a @click="goSearchPage">Search</a>
            <a @click="showMode =0">Edit: {{ fileKey }}</a>
        </n-space>
    </n-space>
</template>