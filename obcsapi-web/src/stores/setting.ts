import { ref, type Ref, onMounted } from 'vue'
import { defineStore } from 'pinia'
import { ObcsapiMentionGet } from "@/api/obcsapi"

export const LocalSetting = defineStore('setting', () => {
    const frontSize = ref("14px");
    const mention: Ref<Array<{label:string,value:string}>> = ref([]);
    const recentEditList: Ref<string[]>= ref([]);
    recentEditList.value = JSON.parse(localStorage.getItem('recentEditList')||'["test.md"]');

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


    return { mention, frontSize,getMention ,recentEditList}
})
