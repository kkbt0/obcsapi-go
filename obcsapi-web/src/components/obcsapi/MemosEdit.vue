<script setup lang="ts">
import { NThing, NInput, NSpace, NButton, NImage, NCheckbox,NImageGroup } from "naive-ui";
import { ObcsapiPostMemos } from "@/api/obcsapi";
import { ref, onMounted } from "vue";
import { memosData } from "@/stores/memos";
import marked from "marked";
import MemosUpload  from "@/components/obcsapi/MemosUpload.vue";

// filekey: string, line: number, oldText: string, newText: string
const props = defineProps<{
    dayKey: any, // dayKey
    line: number, // line
    memosShowText: string, // oldText show fmt
    memosRaw: string, // oldText
}>()

const memos = memosData();
const edit = ref(false);
const inputText = ref(""); // newText
const showUpload = ref(false);

let picList = new Array<string>(); //
let tasksList = new Array<string>(); //
let tasksCheckedList = new Array<boolean>(); //
let nowMd = "";

onMounted(() => {
    inputText.value = props.memosRaw;
})

function moreAction() {
    inputText.value = props.memosRaw;
    edit.value = !edit.value;
}

function saveMemos() {
    ObcsapiPostMemos(props.dayKey, props.line, props.memosRaw, inputText.value).then(data => {
        memos.memosMap.set(data.date, data)
        window.$message.success("Suceess Send");
        edit.value = false;
        showUpload.value = false;
    }).catch(e => {
        console.log(e);
        window.$message.warning("Err Save: " + e);
    });
}

function delMemos() {
    ObcsapiPostMemos(props.dayKey, props.line, props.memosRaw, "").then(data => {
        memos.memosMap.set(data.date, data)
        window.$message.success("Suceess Del");
        edit.value = false;
    }).catch(e => {
        console.log(e);
        window.$message.warning("Err Del: " + e);
    });
}


function markdown(md: string) {
    // 将图片和勾选框单独提取出来 使用 Vue 模板
    picList = new Array<string>();
    tasksList = new Array<string>();
    tasksCheckedList = new Array<boolean>();
    const picRegex = /!\[.*?\]\((.*?)\)/g;
    let match;
    while ((match = picRegex.exec(md))) {
        picList.push(match[1]);
    }
    md = md.replace(picRegex, '')

    const tasksRegex = /- \[[x ]\] .*/gm;
    while ((match = tasksRegex.exec(md))) {
        if (match[0][3] == 'x') {
            tasksCheckedList.push(true);
        } else {
            tasksCheckedList.push(false);
        }
        tasksList.push(match[0].replace(/- \[[x ]\] /, ''));
    }
    md = md.replace(tasksRegex, '')
    nowMd = md;

    return marked(nowMd || '', {
        breaks: true
    })
}

function handleCheckedChange(taskIndex: number) {
    let text = inputText.value;
    const tasksRegex = /- \[[x ]\] .*/gm;
    text = text.replace(tasksRegex, '').trimEnd() + "\n";
    for (let i = 0; i < tasksCheckedList.length; i++) {
        if (tasksCheckedList[i] == false) {
            text += " - [ ] " + tasksList[i] + "\n";
        } else {
            text += " - [x] " + tasksList[i] + "\n";
        }
    }
    inputText.value = text;
    saveMemos();
}

function imgUrlDeal(text:string) {
  inputText.value += `\n${text}\n`;
}

</script>

<template>
    <n-thing @dblclick="moreAction">
        <template #header v-if="memosShowText.slice(2, 7).match(/[0-2][0-9]\:[0-5][0-9]/g)">
            <small>{{ dayKey }} {{ memosShowText.slice(2, 7) }}</small>
        </template>
        <template #header>
            <small>{{ dayKey }}</small>
        </template>

        <template #header-extra><a @click="moreAction">More</a></template>

        <template #description v-if="!edit && memosShowText.slice(2, 7).match(/[0-2][0-9]\:[0-5][0-9]/g)">
            <div v-html="markdown(memosShowText.slice(7))" class="memos"></div>
            <n-space v-if="tasksList.length != 0" vertical>
                <n-checkbox v-for="(task, taskIndex) in tasksList" :key="taskIndex" :label="task"
                    v-model:checked="tasksCheckedList[taskIndex]" @update:checked="handleCheckedChange(taskIndex)" />
            </n-space>
            <n-image-group v-if="picList.length != 0">
                <n-space >
                <n-image v-for="(picUrl,urlIndex) in picList" :key="urlIndex" width="100" :src=picUrl />
            </n-space>
            </n-image-group>

        </template>

        <template #description v-else-if="edit">
            <n-space vertical>
                <n-input v-model:value="inputText" type="textarea" class="memos-input" placeholder="Input Memos"
                    :autosize="{ minRows: 3 }" />
                <n-space justify="space-between">
                    <n-button ghost type="error" @click="delMemos">Del</n-button>
                    <n-button ghost type="info" @click="showUpload=!showUpload">Img</n-button>
                    <n-button ghost type="primary" @click="saveMemos">Save</n-button>
                </n-space>
                <MemosUpload v-if="showUpload" @upload-callback="imgUrlDeal" />
            </n-space>
        </template>

    </n-thing>
</template>

<style scoped>
.memos :deep() img {
    max-width: 40%;
}
</style>