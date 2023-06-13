import { ref } from 'vue'
import { defineStore } from 'pinia'
import { ObcsapiGetMemos } from '@/api/obcsapi';
import { LocalSetting } from "@/stores/setting"

export const memosData = defineStore('memos', () => {
    // ref -> basic ; reactive -> object
    const dayBefore = ref(0); // day before 保证每次加载 从应用上次加载的部分开始
    const memosIndexList = ref(new Array<String>());
    const memosMap = ref(new Map<String, any>());
    if (localStorage.getItem("mainMdList") != undefined) {
        memosMap.value = new Map<String, any>(JSON.parse(localStorage.getItem("mainMdList")||"{}"));
    } 
    
    function addDaily(key: string, value: any) {
        let set = new Set(memosIndexList.value);
        set.add(key);
        memosIndexList.value = Array.from(set);
        // memosIndexList.value.sort();
        setMap(key, value);
    }

    function setMap(key: string, value: any) {
        memosMap.value.set(key,value)
        localStorage.setItem("mainMdList", JSON.stringify([...memosMap.value]));

    }



    // Load Memos
    const loadMemosCount = ref(0); // 取负数 然后 waitMoreMemos()
    async function moreMemos() {
        await ObcsapiGetMemos(dayBefore.value).then(data => {
            console.log(data.date)
            addDaily(data.date, data)
            loadMemosCount.value += data.md_text.length // 统计 memos 数量
            dayBefore.value -= 1;
        }).catch(err => {
            console.log(err)
        })
    }
    
    async function waitMoreMemos(needNum: number) {
        console.log(`Loading more ${needNum} memos`)
        let maxRequest = 5
        let tem = loadMemosCount.value + needNum
        while (loadMemosCount.value < tem && maxRequest > 0) {
            await moreMemos()
            maxRequest -= 1
        }
    }

    waitMoreMemos(LocalSetting().localSetting.LoadMemos||20);

    return { memosIndexList, memosMap , dayBefore,addDaily , setMap , waitMoreMemos}
})
