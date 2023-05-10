<script setup lang="ts">
import MemosConfig from "@/components/obcsapi/MemosConfig.vue";
import { useRouter } from "vue-router";
import { ref,onMounted } from "vue";
import { NTabPane, NTabs, NInputNumber, NSelect } from "naive-ui";
import { LocalSetting } from "@/stores/setting"


const router = useRouter()
const setting = LocalSetting()
const frontSize = ref(14);

frontSize.value = parseInt(setting.frontSize);

function clearCache() {
    localStorage.removeItem("mainMdList")
    window.$message.info("Clearing")
}
function reLogin() {
    localStorage.removeItem("token")
    window.$message.info("注销")
    router.push("/login")
}

const themeMode = ref("跟随系统");

function saveSetting() {
    localStorage.setItem("theme", JSON.stringify({ frontSize: `${frontSize.value}px` }))
    setting.frontSize = `${frontSize.value}px`;
    localStorage.setItem("theme-mode",themeMode.value.toString());
    location.reload();
}

let themeModeOptions = [{ label: "跟随系统", value: "" }, { label: "暗色模式", value: "dark-mode" }, { label: "浅色模式", value: "light-mode" }]

onMounted(() => {
    themeMode.value = localStorage.getItem("theme-mode")||"跟随系统";
})

</script>
<template>
    <n-tabs animated type="card">
        <n-tab-pane name="broSetting" tab="Setting">

            <n-space vertical>
                <div>字体</div>
                <n-input-number v-model:value="frontSize">
                    <template #prefix>字体大小</template>
                    <template #suffix>px</template>
                </n-input-number>
                <n-select v-model:value="themeMode" :options="themeModeOptions" />
                <n-space>
                    <n-button @click="saveSetting" type="info" quaternary>保存设置</n-button>
                    <n-button @click="clearCache" type="info" quaternary>清除缓存</n-button>
                    <n-button @click="reLogin" type="warning" quaternary>注销</n-button>
                </n-space>
            </n-space>


        </n-tab-pane>
        <n-tab-pane name="serverSetting" tab="Server Setting">
            <MemosConfig />
        </n-tab-pane>
    </n-tabs>
</template>