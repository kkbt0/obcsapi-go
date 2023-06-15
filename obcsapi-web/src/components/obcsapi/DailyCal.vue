<script setup lang="ts">
import { ref } from "vue";
import { NCalendar } from "naive-ui";
import { memosData } from "@/stores/memos"
import { LocalSetting } from "@/stores/setting";

const memos = memosData();
const value = ref(Date.now().valueOf());
const emit = defineEmits(["cal-click"]);

function handleUpdateValue(
    _: number,
    { year, month, date }: { year: number; month: number; date: number }
) {
    emit("cal-click", GetFileKey(year, month, date));
}

function GetFileKey(year: number, month: number, date: number): string {
    const f_month = month.toString().padStart(2, '0');
    const f_date = date.toString().padStart(2, '0');
    let template = LocalSetting().localSetting.CalObDaily;
    let reslut = template.replace("2006", "T_Y").replace("01", "T_M").replace("02", "T_D")
    return LocalSetting().localSetting.CalObDailyFir + reslut.replace("T_Y", year.toString()).replace("T_M", f_month.toString()).replace("T_D", f_date.toString())
}

function GetMemosFromDate(year: number, month: number, date: number): string {
    let str = memos.memosMap.get(GetFileKey(year, month, date))
    if (str != undefined) {
        return str.md_text.length
    } else {
        return "";
    }

}
</script>
<template>
    <n-calendar v-model:value="value" #="{ year, month, date }" @update:value="handleUpdateValue">
        {{ GetMemosFromDate(year, month, date) }}
    </n-calendar>
</template>