import { ref, type Ref } from 'vue'
import { defineStore } from 'pinia'

export const TalkStore = defineStore('talk', () => {
    const messages: Ref<string[]> = ref([]);

    return { messages }
})
