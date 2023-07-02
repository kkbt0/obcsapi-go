<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import { NFormItem, NInput, NButton, NButtonGroup } from "naive-ui";
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
    if (data.code == 200 || data.code == 201) {
      localStorage.setItem("token", data.token)
      window.$message.success("登录成功")
      router.push("/").then(() => {
        location.reload();
      })
    } else {
      window.$message.error("登录失败" + JSON.stringify(data))
    }

  }).catch(err => {
    window.$message.error("登录请求失败" + err)
    console.log(err)
  })
}

function LoginByOAuth2() {
  localStorage.setItem("host", formValue.value.host);
  ResetServerHost();
  let url = (localStorage.getItem("host") || "http://localhost:8900") + "/auth/oauth2-login"
  window.location.href = url;
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
  <n-button-group>
    <n-button @click="Login">登录</n-button>
    <n-button @click="LoginByOAuth2">Gitee 登录</n-button>
  </n-button-group>
</template>