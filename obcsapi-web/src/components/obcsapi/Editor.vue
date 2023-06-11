<script setup lang="ts">
import Editor from '@toast-ui/editor';
import '@toast-ui/editor/dist/toastui-editor.css';
import '@toast-ui/editor/dist/theme/toastui-editor-dark.css';
import '@toast-ui/editor/dist/i18n/zh-cn';

import { onMounted, watch } from "vue";
import { ObcsapiTextGet, ObcsapiTextPost } from "@/api/obcsapi"
import { LocalSetting } from "@/stores/setting"
const props = defineProps<{
    fileKey: string, // dayKey
}>()

let darkTheme = true;

if (localStorage.getItem("theme-mode") == "dark-mode") {
    document.documentElement.setAttribute("theme-mode", "dark-mode");
} else if (localStorage.getItem("theme-mode") == "light-mode") {
    document.documentElement.setAttribute("theme-mode", "light-mode");
    darkTheme = false;
} else if (window.matchMedia('(prefers-color-scheme: light)').matches) {
    darkTheme = false;
}


let editor: Editor;
onMounted(() => {
    EditorFontSize();
    EditorInit();
    GetMdText();
})

watch(() => props.fileKey, (newVal) => {
    GetMdText();
    editor.setHeight("85vh");
});

function EditorFontSize() {
    let el = document.getElementById("editor")
    el?.style.setProperty('--theme-font-size', LocalSetting().frontSize);
}

function EditorInit() {
    editor = new Editor({
        el: document.querySelector('#editor'),
        language: 'zh-CN',
        previewStyle: 'vertical',
        initialEditType: 'wysiwyg',
        theme: darkTheme ? 'dark' : 'light',
        setHeight: "85vh",
    });
    editor.insertToolbarItem({ groupIndex: 0, itemIndex: 0 }, {
        el: createLastButton(),
        tooltip: '保存',
        text: '保存',
        className: 'toastui-editor-toolbar-icons first',
        style: { backgroundImage: 'none' }
    });
}

function GetMdText() {
    ObcsapiTextGet(props.fileKey).then(text => {
        editor.setMarkdown(text, true)
        editor.setHeight("85vh");
    }).catch(e => {
        window.$message.warning(e)
    })
}
function createLastButton() {
    const button = document.createElement('button');

    button.className = 'toastui-editor-toolbar-icons first';
    button.style.backgroundImage = 'none';
    button.style.margin = '0';
    button.innerHTML = `💾`;
    button.addEventListener('click', () => {
        SaveMarkdown();
    });
    return button;
}

function SaveMarkdown() {
    let text = editor.getMarkdown().replace(/^(\*\*)/gm, "--").replace(/^(\*)/gm, "-");
    ObcsapiTextPost(props.fileKey, text).then(res => {
        if (res.status == 200) {
            window.$message.success("Success!");
            SaveRecentEditList(props.fileKey)
        } else {
            window.$message.error("错误,未完成保存")
        }
    }).catch(e => {
        window.$message.warning(e)
    })
}

function SaveRecentEditList(fileKey: string) {
    // 去重
    const index = LocalSetting().recentEditList.indexOf(fileKey);

    if (index > -1) {
        LocalSetting().recentEditList.splice(index, 1)
    }

    // 新增
    LocalSetting().recentEditList.unshift(fileKey);
    if (LocalSetting().recentEditList.length > 10) {
        LocalSetting().recentEditList.splice(10)
    }
    localStorage.setItem('recentEditList', JSON.stringify(LocalSetting().recentEditList));
}

</script>

<template>
    <div id="editor" class="editor"></div>
</template>
<style scoped>
.editor :deep() .toastui-editor-contents {
    font-size: var(--theme-font-size);
}

.editor :deep() .ProseMirror {
    font-size: var(--theme-font-size)
}
</style>