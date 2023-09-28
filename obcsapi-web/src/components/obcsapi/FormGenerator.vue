<script setup lang="ts">
import { ref } from "vue"
import VueForm from "@lljj/vue3-form-naive"
import { ObcsapiFormPost } from "@/api/obcsapi";
// https://1.xrender.fun/generator
// https://form.lljj.me/v3/#/demo?type=Simple&ui=VueNaiveForm

const formData = ref();
const formJsonShemeText = ref("");
const formJsonSheme = ref({
    "type": "object",
    "labelWidth": 120,
    "displayType": "row",
    "properties": {
        "text1": {
            "title": "单行1 text1",
            "type": "string",
            "props": {}
        },
        "text2": {
            "title": "单行2 text1",
            "type": "string",
            "props": {}
        }
    }
})
const result = ref("");
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
function updateFormJsonSheme() {
    try {
        formJsonSheme.value = JSON.parse(formJsonShemeText.value)
    } catch (error) {
        window.$message.error(`${error}`)
    }
}

</script>
<template>
    <n-input v-model:value="formJsonShemeText" type="textarea" placeholder="JsonSheme" />
    <n-button @click="updateFormJsonSheme">updateFormJsonSheme</n-button>
    <vue-form v-model="formData" :schema="formJsonSheme" @cancel="handlerCancel" @submit="handlerSubmit" />
    <div>{{ result }}</div>
</template>

<style scoped></style>
