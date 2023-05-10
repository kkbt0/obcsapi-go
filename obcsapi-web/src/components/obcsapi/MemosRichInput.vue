
  
<script lang="ts">
import marked from "marked";
import TurndownService from 'turndown'
import { NSpace, NButton, useThemeVars } from "naive-ui";

export default {

    props: ['value'],
    data() {
        return {
            innerValue: marked(this.value || '', {
                breaks: true
            })
        }
    },

    mounted() {
        document.execCommand('defaultParagraphSeparator', false, 'p')
    },

    methods: {
        onInput(event: any) {
            console.log(event.target.innerHTML);
            const turndown = new TurndownService({
                emDelimiter: '_',
                linkStyle: 'inlined',
                headingStyle: 'atx'
            })
            console.log(turndown.turndown(event.target.innerHTML))
            this.$emit('input', turndown.turndown(event.target.innerHTML))
        },
        applyBold() {
            document.execCommand('bold')
        },
        applyItalic() {
            document.execCommand('italic')
        },
        applyHeading() {
            document.execCommand('formatBlock', false, '<h1>')
        },
        applyUl() {
            document.execCommand('insertUnorderedList')
        },
        applyOl() {
            document.execCommand('insertOrderedList')
        },
        undo() {
            document.execCommand('undo')
        },
        redo() {
            document.execCommand('redo')
        }
    }
}
</script>
<template>
    <div>
        <n-space>
            <n-button quaternary type="primary" @click="applyBold">B</n-button>
            <n-button quaternary type="primary" @click="applyItalic">I</n-button>
            <n-button quaternary type="primary" @click="applyHeading">H</n-button>
            <n-button quaternary type="primary" @click="applyUl">·</n-button>
            <n-button quaternary type="primary" @click="applyOl">1.</n-button>
            <n-button quaternary type="primary" @click="undo">&lt;-</n-button>
            <n-button quaternary type="primary" @click="redo">-&gt;</n-button>
        </n-space>
        {{ value }}
        <div @input="onInput" v-html="innerValue" contenteditable="true" class="memosRichInput" />
    </div>
</template>

<style scoped>
.memosRichInput {
    background-color: rgba(255, 255, 255, 0.1);
    border-color: #63e2b7;
    outline: none;
}
</style>