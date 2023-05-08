import { ref } from 'vue'
import { defineStore } from 'pinia'

export const memosData = defineStore('memos', () => {
    // ref -> basic ; reactive -> object
    const memosIndexList = ref(new Array<String>());
    const memosMap = ref(new Map<String, any>());
    function addDaily(key: string, value: any) {
        let set = new Set(memosIndexList.value);
        set.add(key);
        memosIndexList.value = Array.from(set);
        memosIndexList.value.sort();
        memosMap.value.set(key, value);
    }
    return { memosIndexList, memosMap ,addDaily}
})
