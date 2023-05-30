import { ref, type Ref, onMounted } from 'vue'
import { defineStore } from 'pinia'
import { ObcsapiMentionGet } from "@/api/obcsapi"

export const LocalSetting = defineStore('setting', () => {
    const frontSize = ref("14px");
    const mention: Ref<Array<string>> = ref([]);

    onMounted(() => {
        ObcsapiMentionGet().then(obj => {
            if(obj.tags!=null) {
                mention.value = obj.tags;
            }
        });
        
    })

    frontSize.value = JSON.parse(localStorage.getItem("theme") || "{}").frontSize


    return { mention, frontSize }
})
