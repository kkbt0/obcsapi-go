<script setup lang="ts">
import { ref, watch, type Ref } from "vue"
import VueForm from "@lljj/vue3-form-naive"
import { ObcsapiFormPost } from "@/api/obcsapi";
import { LocalSetting } from "@/stores/setting"

// https://1.xrender.fun/generator
// https://form.lljj.me/v3/#/demo?type=Simple&ui=VueNaiveForm

const mentionList: Ref<Array<any>> = ref([]);
const formData = ref();
const formJsonShemeText: Ref<string> = ref(`{}`);
const formJsonSheme:Ref<any> = ref();

const result = ref("");

selectOptionsInit();

function handlerSubmit() {
    console.log(formData.value)
    ObcsapiFormPost(formData.value)
        .then((response) => response.text())
        .then((data) => {
            console.log(data)
            result.value = data;
            window.$message.info(data);
        });
}
function handlerCancel() {
    console.log("cancel")
}

function selectOptionsInit() {
    let options = [];
    let rawOptions = LocalSetting().formJSONSchemaOptions;
    for (let i = 0; i < rawOptions.length; i++) {
        options.push({
            label: rawOptions[i].title,
            value: rawOptions[i].json_schema
        })
    }
    mentionList.value = options;
    if(mentionList.value.length > 0){
        formJsonShemeText.value = mentionList.value[0].value;
    }
    formJsonSheme.value = JSON.parse(formJsonShemeText.value);
}

watch(formJsonShemeText, () => {
    console.log(formJsonShemeText.value)
    console.log(JSON.stringify(formJsonSheme.value))
    formJsonSheme.value = JSON.parse(formJsonShemeText.value);
    console.log(JSON.stringify(formJsonSheme.value))
    // bug
})
 
</script>
<template>
    <n-space vertical>
        {{  formJsonShemeText }}
        <n-select v-model:value="formJsonShemeText" :options="mentionList" />
        <vue-form v-model="formData" :schema="formJsonSheme" @cancel="handlerCancel" @submit="handlerSubmit" />
        <div v-text="result" style="white-space: pre-wrap;"></div>
    </n-space>
</template>

<style scoped></style>
