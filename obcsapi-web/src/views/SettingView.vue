<script setup lang="ts">
import MemosConfig from "@/components/obcsapi/MemosConfig.vue";
import { useRouter } from "vue-router";
import { ref, onMounted, type Ref } from "vue";
import { NTabPane, NTabs, NInputNumber, NSelect, NDynamicInput } from "naive-ui";
import { LocalSetting } from "@/stores/setting"
import { ObcsapiConfigPost } from "@/api/obcsapi"

const router = useRouter()
const setting = LocalSetting()
const frontSize = ref(14);
const mentionList: Ref<Array<string>> = ref([]);

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

function saveMention() {
    ObcsapiConfigPost({ "mention": { "tags": mentionList.value } }).then(text => {
        if (text == "Success") {
            window.$message.success("保存成功")
            LocalSetting().getMention()
        }
    })
}

const themeMode = ref("跟随系统");

function saveSetting() {
    localStorage.setItem("theme", JSON.stringify({ frontSize: `${frontSize.value}px` }))
    setting.frontSize = `${frontSize.value}px`;
    localStorage.setItem("theme-mode", themeMode.value.toString());
    localStorage.setItem("LocalSetting",JSON.stringify(LocalSetting().localSetting));
    location.reload();
}

const themeModeOptions = [{ label: "跟随系统", value: "" }, { label: "暗色模式", value: "dark-mode" }, { label: "浅色模式", value: "light-mode" }]

onMounted(() => {
    themeMode.value = localStorage.getItem("theme-mode") || "跟随系统";
    console.log(LocalSetting().localSetting.LoadMemos);
    getMention(); // 初始化这个组件的列表
})

function getMention() {
    LocalSetting().mention.forEach(obj => {
        mentionList.value.push(obj.value);
    })
}

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
                <div>主题</div>
                <n-select v-model:value="themeMode" :options="themeModeOptions" />
                <div>自动对焦</div>
                <n-switch v-model:value="LocalSetting().localSetting.AutoFocus" />
                <div>每次最少加载 Memos 数量 默认20 ；并且每次加载最多请求 5 个文件</div>
                <n-input-number v-model:value="LocalSetting().localSetting.LoadMemos">
                    <template #suffix>条</template>
                </n-input-number>
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
        <n-tab-pane name="serverMention" tab="Mention">
            <n-space vertical>
                <a>提示词输入框 # 触发</a>

                <n-dynamic-input v-model:value="mentionList" placeholder="请输入提示词" :min="0" />
                <n-space justify="end">
                    <n-button @click="saveMention" type="info" quaternary>保存提示词</n-button>
                </n-space>
            </n-space>

        </n-tab-pane>
    </n-tabs>
</template>