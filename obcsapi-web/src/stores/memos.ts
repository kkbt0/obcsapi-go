import { ref } from 'vue'
import { defineStore } from 'pinia'

export const memosData = defineStore('memos', () => {
    // ref -> basic ; reactive -> object
    const memosIndexList = ref(new Array<String>());
    const memosMap = ref(new Map<String, any>());
    if (localStorage.getItem("mainMdList") != undefined) {
        memosMap.value = new Map<String, any>(JSON.parse(localStorage.getItem("mainMdList")||"{}"));
    } 
    
    function addDaily(key: string, value: any) {
        let set = new Set(memosIndexList.value);
        set.add(key);
        memosIndexList.value = Array.from(set);
        memosIndexList.value.sort();
        setMap(key, value);
    }

    function setMap(key: string, value: any) {
        memosMap.value.set(key,value)
        localStorage.setItem("mainMdList", JSON.stringify([...memosMap.value]));

    }
    return { memosIndexList, memosMap ,addDaily , setMap}
})
