<script setup lang="ts">
import { NThing, NSpace, NButton, NImage, NCheckbox, NImageGroup, NDropdown, NScrollbar, NMention } from "naive-ui";
import { ObcsapiPostMemos } from "@/api/obcsapi";
import { ref, onUpdated, onMounted, watch } from "vue";
import { memosData } from "@/stores/memos";
import marked from "marked";
import MemosUpload from "@/components/obcsapi/MemosUpload.vue";
import { LocalSetting } from "@/stores/setting"

// filekey: string, line: number, oldText: string, newText: string
const props = defineProps<{
    dayKey: any, // dayKey
    line: number, // line
    memosShowText: string, // oldText show fmt
    memosRaw: string, // oldText
}>()

const memos = memosData();
const edit = ref(false);
const inputText = ref(props.memosRaw); // newText
const showUpload = ref(false);

let picList = new Array<string>(); //
let tasksList = new Array<string>(); //
let tasksCheckedList = new Array<boolean>(); //
let nowMd = "";

onUpdated(() => {
    inputText.value = props.memosRaw;
})

watch(() => inputText.value, (newVal) => {
    LocalSetting().lastInputPush(newVal);
});

function moreAction() {
    edit.value = !edit.value;
}

function saveMemos() {
    ObcsapiPostMemos(props.dayKey, props.line, props.memosRaw, inputText.value).then(data => {
        if (data.md_text != undefined) {
            memos.setMap(data.date, data)
            window.$message.success("Suceess Send");
            edit.value = false;
            showUpload.value = false;
            localStorage.setItem("lastInput", "");
        } else {
            window.$message.error("Err Save: " + JSON.stringify(data));
        }

    }).catch(e => {
        console.log(e);
        window.$message.warning("Err Save: " + e);
    });
}

function delMemos() {
    LocalSetting().delMemosListPush(inputText.value);
    ObcsapiPostMemos(props.dayKey, props.line, props.memosRaw, "").then(data => {
        if (data.md_text != undefined) {
            memos.setMap(data.date, data)
            window.$message.success("Suceess Del");
            edit.value = false;
        } else {
            window.$message.error("Err Save: " + JSON.stringify(data));
        }

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
    if (isCollapsed.value) {
        isCollapsed.value = false;
        return;
    }
    let text = inputText.value;
    let isEmpty = false;
    const tasksRegex = /- \[[x ]\] .*/gm;
    text = text.replace(tasksRegex, '').trimEnd() + "\n";
    console.log(text);
    if (text.trim() == "") { // 被去除的什么都没有了 特别是之只有一个选项撑起来的情况
        isEmpty = true;
    }
    for (let i = 0; i < tasksCheckedList.length; i++) {
        if (tasksCheckedList[i] == false) {
            text += " - [ ] " + tasksList[i] + "\n";
        } else {
            text += " - [x] " + tasksList[i] + "\n";
        }
    }
    if (isEmpty) {
        inputText.value = text.slice(2);// 把第一个空格去掉，以撑起来一个 memos
    } else {
        inputText.value = text;
    }
    saveMemos();
}

function imgUrlDeal(text: string) {
    inputText.value += `\n${text}\n`;
}

const options = ref([
    {
        label: '修改',
        key: 0
    },
    {
        label: '移动今日',
        key: 1
    },
    {
        label: '转为 TODO',
        key: 2
    },
    {
        label: '删除',
        key: 3
    }
])

function handleSelect(key: string | number) {
    if (key == 0) {
        edit.value = !edit.value;
    } else if (key == 1) {
        let realText = ""
        if (inputText.value.slice(2, 7).match(/[0-2][0-9]\:[0-5][0-9]/g)) {
            realText = inputText.value.slice(7)
        } else {
            realText = inputText.value.slice(2)
        }
        console.log(realText)
        LocalSetting().delMemosListPush(inputText.value);
        ObcsapiPostMemos(props.dayKey, props.line, props.memosRaw, "").then(data => { // 先删除
            if (data.md_text != undefined) {
                memos.setMap(data.date, data)
                window.$message.success("Suceess Del");
                edit.value = false;
                ObcsapiPostMemos("", 9999, "", realText).then(data => { //再新增
                    memos.setMap(data.date, data)
                    window.$message.success("Suceess ReSend");
                    edit.value = false;
                    showUpload.value = false;
                }).catch(e => {
                    console.log(e);
                    window.$message.warning("Err ReSend: " + e);
                });
            } else {
                window.$message.error("Err Save: " + JSON.stringify(data));
            }

        }).catch(e => {
            console.log(e);
            window.$message.warning("Err Del: " + e);
        });
    } else if (key == 2) { // - xxx -> - [ ] xx
        let arr = inputText.value.split("\n");
        for (var i = 0; i < arr.length; i++) {
            if (arr[i].trim().startsWith("- [")) { // do nothing
            } else if (arr[i].trim().startsWith("- ")) {
                arr[i] = "- [ ]" + arr[i].trim().slice(1)
            } else { // xx
                arr[i] = " - [ ] " + arr[i].trim()
            }
        }
        inputText.value = arr.join("\n");
        saveMemos();
    }
    else if (key == 3) {
        delMemos()
    }
}
const isCollapsed = ref(true);
const checkNeedCollapse = ref(false);

// 截取几行文字，返回的字符串最少100字符
function truncateText(text: string, lines: number): string {
    // 将文本按行分割成数组
    const textLines = text.split('\n');
    const truncatedLines = textLines.slice(0, lines);
    // 将截取的行重新组合成一个字符串
    const truncatedText = truncatedLines.join('\n');
    // 如果截取后的文本长度小于100个字符，则继续截取更多行
    if (truncatedText.length < 100 && lines < textLines.length) {
        return truncateText(text, lines + 1);
    }
    if (truncatedText == text) {
        return truncatedText;
    } else {
        return truncatedText + '……';
    }
}

onMounted(() => {
    checkNeedCollapse.value = truncateText(props.memosShowText, 1) != props.memosShowText;
    isCollapsed.value = checkNeedCollapse.value;
})

</script>

<template>
    <n-thing @dblclick="moreAction">
        <template #header v-if="memosShowText.slice(2, 7).match(/[0-2][0-9]\:[0-5][0-9]/g)">
            <small>{{ dayKey }} <div style="font-weight: bold;display:inline;">{{ memosShowText.slice(2, 7) }}</div></small>
        </template>
        <template #header>
            <small>{{ dayKey }}</small>
        </template>

        <template #header-extra>
            <n-space>
                <n-button v-if="checkNeedCollapse" @click="isCollapsed = !isCollapsed;" quaternary type="primary">{{
                    isCollapsed ? '◀' : '▼' }}</n-button>
                <n-dropdown trigger="hover" :options="options" @select="handleSelect">
                    <n-button quaternary>···</n-button>
                </n-dropdown>
            </n-space>
        </template>
        <!-- - 12:34 xxx -->
        <template #description v-if="!edit && memosShowText.slice(2, 7).match(/[0-2][0-9]\:[0-5][0-9]/g)">
            <n-scrollbar x-scrollable class="scrollbar-content">
                <div v-if="!isCollapsed" v-html="markdown(memosShowText.slice(7))" class="memos"></div>
                <div v-else v-html="markdown(truncateText(memosShowText.slice(7), 1))" class="memos"></div>
                <n-space v-if="tasksList.length != 0" vertical>
                    <n-checkbox v-for="(task, taskIndex) in tasksList" :key="taskIndex" :label="task"
                        v-model:checked="tasksCheckedList[taskIndex]" @update:checked="handleCheckedChange(taskIndex)" />
                </n-space>
                <n-image-group v-if="picList.length != 0">
                    <n-space>
                        <n-image v-for="(picUrl, urlIndex) in picList" :key="urlIndex" width="100" loading="lazy"
                            :src=picUrl />
                    </n-space>
                </n-image-group>
            </n-scrollbar>
        </template>
        <!-- - xxx -->
        <template #description v-else-if="!edit">
            <n-scrollbar x-scrollable class="scrollbar-content">
                <div v-if="!isCollapsed" v-html="markdown(memosShowText)" class="memos"></div>
                <div v-else v-html="markdown((truncateText(memosShowText.slice(0),1)))" class="memos"></div>
                <n-space v-if="tasksList.length != 0" vertical>
                    <n-checkbox v-for="(task, taskIndex) in tasksList" :key="taskIndex" :label="task"
                        v-model:checked="tasksCheckedList[taskIndex]" @update:checked="handleCheckedChange(taskIndex)" />
                </n-space>
                <n-image-group v-if="picList.length != 0">
                    <n-space>
                        <n-image v-for="(picUrl, urlIndex) in picList" :key="urlIndex" width="100" :src=picUrl />
                    </n-space>
                </n-image-group>
            </n-scrollbar>
        </template>

        <template #description v-if="edit">
            <n-space vertical>
                <n-mention v-model:value="inputText" type="textarea" class="memos-input" placeholder="Memos"
                    :autosize="{ minRows: 3 }" :options="LocalSetting().mention" :prefix="['#']" />
                <n-space justify="space-between">
                    <n-button quaternary type="error" @click="delMemos">❌</n-button>
                    <n-button quaternary type="info" @click="showUpload = !showUpload">📌</n-button>
                    <n-button quaternary @click="edit = !edit">🔙</n-button>
                    <n-button quaternary type="primary" @click="saveMemos">💾</n-button>
                </n-space>
                <MemosUpload v-if="showUpload" @upload-callback="imgUrlDeal" />
            </n-space>
        </template>

        <template #action v-if="checkNeedCollapse && !isCollapsed">
            <n-space justify="end">
                <n-button @click="isCollapsed = !isCollapsed;" quaternary type="primary">{{ isCollapsed ? '展开' : '折叠'
                }}</n-button>
            </n-space>
        </template>
    </n-thing>
</template>

<style scoped>
.memos :deep() img {
    max-width: 40%;
    overflow: hidden;
    transition: max-height 0.3s ease;
    max-height: fit-content;
}

.memos :deep() pre {
    max-width: 80vw;
}

.memos :deep() p {
    word-break: break-all;
}
</style>