import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

export const LocalSetting = defineStore('setting', () => {
    const count = ref(0)
    const frontSize = ref("14px");

    frontSize.value = JSON.parse(localStorage.getItem("theme") || "{}").frontSize

    const doubleCount = computed(() => {
        console.log("double count")
        count.value * 2
    })


    return { count, doubleCount ,frontSize }
})
