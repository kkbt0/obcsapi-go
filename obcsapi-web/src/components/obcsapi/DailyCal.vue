<script setup lang="ts">
import { ref } from "vue";
import { isYesterday, addDays } from 'date-fns/esm';
import { NCalendar } from "naive-ui";
import { memosData } from "@/stores/memos"

const memos = memosData();
const value = ref(addDays(Date.now(), 1).valueOf());
const emit = defineEmits(["cal-click"]);

function handleUpdateValue(
    _: number,
    { year, month, date }: { year: number; month: number; date: number }
) {
    emit("cal-click",GetFileKey(year, month, date));
}

function GetFileKey(year: number,month: number,date: number) :string {
    return `${year}-${month.toString().padStart(2,'0')}-${date.toString().padStart(2,'0')}`
}

function GetMemosFromDate(year: number,month: number,date: number): string {
    let str = memos.memosMap.get(GetFileKey(year,month,date))
    if (str != undefined ){
        return  str.md_text.length
    } else {
        return "";
    }

}
</script>
<template>
    <n-calendar v-model:value="value" #="{ year, month, date }" 
        @update:value="handleUpdateValue">
        {{  GetMemosFromDate(year,month,date) }}
        <!-- {{ year }}-{{ month }}-{{ date }} -->
    </n-calendar>
</template>