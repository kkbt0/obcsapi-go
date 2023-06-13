import { ref, type Ref, onMounted, watch } from 'vue'
import { defineStore } from 'pinia'
import { ObcsapiMentionGet } from "@/api/obcsapi"

class LocalSettingsClass {
    AutoFocus: boolean = false
    LoadMemos: number = 20
}

export const LocalSetting = defineStore('setting', () => {
    const frontSize = ref("14px");
    const mention: Ref<Array<{label:string,value:string}>> = ref([]);
    const recentEditList: Ref<string[]>= ref([]);
    const allFileKeyList: Ref<string[]>= ref([]);
    const localSetting: Ref<LocalSettingsClass> = ref(JSON.parse(localStorage.getItem('LocalSetting')||'{}'));
    recentEditList.value = JSON.parse(localStorage.getItem('recentEditList')||'["test.md"]');
    allFileKeyList.value = JSON.parse(localStorage.getItem('AllFileKeyList')||'["test.md"]');

    onMounted(() => {
        getMention()
    })

    function getMention() {
        ObcsapiMentionGet().then(obj => {
            console.log("Load Mention")
            if (obj.tags != null) {
                mention.value = []
                obj.tags.forEach((val: string) => {
                    mention.value.push({ label: val, value: val });
                })
            }
        });
    }

    frontSize.value = JSON.parse(localStorage.getItem("theme") || "{}").frontSize


    return { mention, frontSize,getMention ,recentEditList , allFileKeyList , localSetting}
})
