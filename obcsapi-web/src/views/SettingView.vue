<script setup lang="ts">
import MemosConfig from "@/components/obcsapi/MemosConfig.vue";
import { useRouter } from "vue-router";
import { ref, onMounted, type Ref } from "vue";
import { NTabPane, NTabs, NInputNumber, NSelect, NDynamicInput, NList, NListItem,NCollapse,NCollapseItem } from "naive-ui";
import { LocalSetting } from "@/stores/setting"
import { memosData } from "@/stores/memos";
import { ObcsapiConfigPost, ObcsapiUpdateCache } from "@/api/obcsapi"

const router = useRouter()
const setting = LocalSetting()
const frontSize = ref(14);
const mentionList: Ref<Array<string>> = ref([]);
const updateFileKey = ref("");

frontSize.value = parseInt(setting.frontSize);

function clearCache() {
    localStorage.removeItem("mainMdList")
    localStorage.removeItem("mainMdListIndex")
    localStorage.removeItem("delMemosList")
    localStorage.removeItem("AllFileKeyList")
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
    localStorage.setItem("LocalSetting", JSON.stringify(LocalSetting().localSetting));
    location.reload();
}

const themeModeOptions = [{ label: "跟随系统", value: "" }, { label: "暗色模式", value: "dark-mode" }, { label: "浅色模式", value: "light-mode" }]

onMounted(() => {
    themeMode.value = localStorage.getItem("theme-mode") || "跟随系统";
    getMention(); // 初始化这个组件的列表
})

function getMention() {
    LocalSetting().mention.forEach(obj => {
        mentionList.value.push(obj.value);
    })
}

function updateServerCache() {
    if(updateFileKey.value!="") {
        ObcsapiUpdateCache(updateFileKey.value).then( res => {
            if(res.status==200) {
                window.$message.success("Success")
            } else {
                window.$message.warning(`${res.status}`)
            }
        }).catch( err => {
            window.$message.error(err)
        });
    } else {
        window.$message.warning("Empty")
    }
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
                <div>每次最少加载 Memos 数量 默认20 ；并且每次加载最多请求 5 个文件</div>
                <n-input-number v-model:value="LocalSetting().localSetting.LoadMemos">
                    <template #suffix>条</template>
                </n-input-number>
                <div>加载前使用浏览器缓存</div>
                <n-switch v-model:value="LocalSetting().localSetting.UseCacheFirst" />
                <div>加载前使用浏览器缓存文件个数，如果数量过大初次渲染时间会较长</div>
                <n-input-number v-model:value="LocalSetting().localSetting.UseCacheFileNum">
                    <template #suffix>个文件</template>
                </n-input-number>
                <div>已缓存文件列表 {{ LocalSetting().allFileKeyList.length }} 项</div>
                <div>存在的日记文件索引 {{ memosData().memosIndexList.length }} 个</div>
                <div>已缓存日记文件 {{ memosData().memosMap.size }} 个</div>
                <div>已删除 {{ LocalSetting().delMemosList.length }} 个 Memos </div>
                <n-collapse>
                    <n-collapse-item title="已删除 Memos" name="1">
                        <n-list v-for="(val, index) in LocalSetting().delMemosList" :key="index">
                            <n-list-item>{{ val }}</n-list-item>
                        </n-list>
                    </n-collapse-item>
                    <n-collapse-item title="最后输入缓存" name="2">
                        {{   LocalSetting().lastInput  }}
                    </n-collapse-item>
                    <n-collapse-item title="更新服务器指定文件缓存" name="3">
                        <n-input v-model:value="updateFileKey" />
                <n-button @click="updateServerCache"  type="info" quaternary>更新服务器指定文件缓存</n-button>
                    </n-collapse-item>
                </n-collapse>
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