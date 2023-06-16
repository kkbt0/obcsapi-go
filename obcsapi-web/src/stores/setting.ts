import { ref, type Ref, onMounted } from 'vue'
import { defineStore } from 'pinia'
import { ObcsapiConfigGet } from "@/api/obcsapi"

class LocalSettingsClass {
    LoadMemos: number = 20
    UseCacheFirst: boolean = false
    UseCacheFileNum: number = 5
    CalObDailyFir: string = "日记/"
    CalObDaily: string = "2006-01-02"
}

export const LocalSetting = defineStore('setting', () => {
    const webDesc = "v20230616-1855 for server_v4.2.4";
    const frontSize = ref("14px");
    const mention: Ref<Array<{ label: string, value: string }>> = ref([]);
    const recentEditList: Ref<string[]> = ref([]);
    const allFileKeyList: Ref<string[]> = ref([]);
    const lastInput: Ref<string> = ref("");
    const delMemosList: Ref<string[]> = ref([]);
    let timer: string | number | NodeJS.Timeout | undefined = undefined;
    const localSetting: Ref<LocalSettingsClass> = ref(JSON.parse(localStorage.getItem('LocalSetting') || '{}'));
    recentEditList.value = JSON.parse(localStorage.getItem('recentEditList') || '["test.md"]');
    allFileKeyList.value = JSON.parse(localStorage.getItem('AllFileKeyList') || '["test.md"]');
    lastInput.value = localStorage.getItem('lastInput') || "";
    delMemosList.value = JSON.parse(localStorage.getItem('delMemosList') || '[]');

    onMounted(() => {
        getFromServerRunConfig()
    })

    function getFromServerRunConfig() {
        ObcsapiConfigGet().then((config: any) => {
            if (config.mention.tags != null) {
                mention.value = []
                config.mention.tags.forEach((val: string) => {
                    mention.value.push({ label: val, value: val });
                })
            }
            localSetting.value.CalObDailyFir = config.ob_daily_config.ob_daily_dir;
            localSetting.value.CalObDaily = config.ob_daily_config.ob_daily;
        })
    }

    frontSize.value = JSON.parse(localStorage.getItem("theme") || "{}").frontSize

    function lastInputPush(text: string) {
        clearTimeout(timer);
        timer = setTimeout(function () {
            LocalSetting().lastInput = text;
            localStorage.setItem("lastInput", text);
        }, 1000);
    }

    function delMemosListPush(text: string) {
        delMemosList.value.unshift(text);
        localStorage.setItem('delMemosList', JSON.stringify(delMemosList.value));
    }


    return {
        mention, frontSize, recentEditList,getFromServerRunConfig, webDesc,
        allFileKeyList, localSetting,
        lastInput, lastInputPush,
        delMemosListPush, delMemosList
    }
})
