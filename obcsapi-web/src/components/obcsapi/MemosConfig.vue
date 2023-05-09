<script setup lang="ts">
import { ref } from "vue";
import VueForm from "@lljj/vue3-form-naive"
import { NForm,NFormItem,NInput,NInputNumber,NButton } from "naive-ui";

const formData = ref({});
const schema = ref({
    type: 'object',
    required: [
        'userName',
        'age',
    ],
    properties: {
        userName: {
            type: 'string',
            title: '用户名',
            default: 'Liu.Jun',
        },
        age: {
            type: 'number',
            title: '年龄'
        },
        bio: {
            type: 'string',
            title: '签名',
            minLength: 10,
            default: '知道的越多、就知道的越少',
        }
    }
});
function handlerSubmit() {
    console.log(formData.value);
    console.log("Client")
}
</script>
<template>
    <h1>xxx</h1>
    <vue-form v-model="formData" :schema="schema" @submit="handlerSubmit">
    </vue-form>
</template>