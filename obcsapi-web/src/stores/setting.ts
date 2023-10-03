import { ref, type Ref, onMounted } from 'vue'
import { defineStore } from 'pinia'
import { ObcsapiConfigGet } from "@/api/obcsapi"

class LocalSettingsClass {
    Theme: string = ""
    FrontSize: string = "14px"
    LoadMemos: number = 20
    UseCacheFirst: boolean = false
    UseCacheFileNum: number = 5
    CalObDailyFir: string = "日记/"
    CalObDaily: string = "2006-01-02"
}

export const LocalSetting = defineStore('setting', () => {
    const webDesc = "v20231003-1800 for server_v4.2.8";
    const mention: Ref<Array<{ label: string, value: string }>> = ref([]);
    const recentEditList: Ref<string[]> = ref([]);
    const allFileKeyList: Ref<string[]> = ref([]);
    const lastInput: Ref<string> = ref("");
    const delMemosList: Ref<string[]> = ref([]);
    const formJSONSchemaOptions: Ref<any> = ref([]);
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
            if (config.mention.from_options != null && config.mention.from_options.length != 0) {
                formJSONSchemaOptions.value = config.mention.from_options;
            } else {
                formJSONSchemaOptions.value = [{
                    title: 'Demo', json_schema: `{
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
                }`}
                ]
            }
            localSetting.value.CalObDailyFir = config.ob_daily_config.ob_daily_dir;
            localSetting.value.CalObDaily = config.ob_daily_config.ob_daily;
        })
    }


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
        mention, recentEditList, getFromServerRunConfig, webDesc,
        allFileKeyList, localSetting,
        lastInput, lastInputPush,
        delMemosListPush, delMemosList, formJSONSchemaOptions
    }
})
