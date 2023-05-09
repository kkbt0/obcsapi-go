<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import { NFormItem, NInput, NButton } from "naive-ui";
import { ObcsapiLogin, ResetServerHost } from "@/api/obcsapi";

const formValue = ref({ "host": "http://localhost:8900", "username": "", "password": "" });

const router = useRouter()

onMounted(() => {
  formValue.value.host = localStorage.getItem("host") || window.location.protocol + "//" + window.location.host;
})

function Login() {
  localStorage.setItem("host", formValue.value.host);
  ResetServerHost();
  ObcsapiLogin(formValue.value.username, formValue.value.password).then(data => {
    localStorage.setItem("token", data.token)
    window.$message.success("登录成功")
    router.push("/")
  }).catch(err => {
    window.$message.error("登录失败" + err)
    console.log(err)
  })
}


</script>
<template>
  <n-form-item label="服务器">
    <n-input v-model:value="formValue.host" placeholder="服务器" />
  </n-form-item>
  <n-form-item label="用户名">
    <n-input v-model:value="formValue.username" placeholder="用户名" />
  </n-form-item>
  <n-form-item label="密码">
    <n-input v-model:value="formValue.password" type="password" @keydown.enter.prevent />
  </n-form-item>
  <n-button @click="Login">登录</n-button>
</template>