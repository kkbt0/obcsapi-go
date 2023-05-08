<script setup lang="ts">
import { NThing, NInput, NSpace, NButton } from "naive-ui";
import { ObcsapiPostMemos } from "@/api/obcsapi";
import { ref, onMounted } from "vue";
import { memosData } from "@/stores/memos";

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

</script>

<template>
    <!-- 三种情况 --->
    <n-thing v-if="memosShowText.slice(2, 7).match(/[0-2][0-9]\:[0-5][0-9]/g)" @dblclick="moreAction">
        <template #header>
            <small>{{ dayKey }} {{ memosShowText.slice(2, 7) }}</small>
        </template>
        <template #header-extra><a @click="moreAction">More</a></template>
        <template #description v-if="!edit">{{ memosShowText.slice(7) }}</template>
        <template #description v-else>
            <n-space vertical>
                <n-input v-model:value="inputText" type="textarea" class="memos-input" placeholder="Input Memos"
                    :autosize="{ minRows: 3 }" />
                <n-space justify="space-between">
                    <n-button ghost type="error" @click="delMemos">Del</n-button>
                    <n-button ghost type="primary" @click="saveMemos">Save</n-button>
                </n-space>
            </n-space>
        </template>
    </n-thing>
    <n-thing v-else-if="memosShowText.trim()!=''" @dblclick="moreAction">
        <template #header>
            <small>{{ dayKey }}</small>
        </template>
        <template #header-extra><a @click="moreAction">More</a></template>
        <template #description v-if="!edit">{{ memosShowText }}</template>
        <template #description v-else>
            <n-space vertical>
                <n-input v-model:value="inputText" type="textarea" class="memos-input" placeholder="Input Memos"
                    :autosize="{ minRows: 3 }" />
                <n-space justify="space-between">
                    <n-button ghost type="error" @click="delMemos">Del</n-button>
                    <n-button ghost type="primary" @click="saveMemos">Save</n-button>
                </n-space>
            </n-space>
        </template>
    </n-thing>
</template>

<style scoped></style>