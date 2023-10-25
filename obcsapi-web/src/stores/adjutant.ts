import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

export const Adjutant = defineStore('adjutant', () => {
    const showText = ref("Weclome!");
    const textList = ref(["Welcome!"]);
    const nowPoint = ref(0);
    
    let timer = setTimeout(function() {
        showText.value = "Weclome!";
    }, 3000);
    // 显示 3s
    function setText(text: string) {
        textList.value.push(text);
        nowPoint.value++;
        if (nowPoint.value >= textList.value.length) {
            nowPoint.value = textList.value.length - 1;
        }
        const nowText = textList.value[nowPoint.value];
        // showText.value = nowText;  // 头部显示
        window.$message.success(nowText); // 弹窗显示
        clearTimeout(timer); // 3s 重新计数
        timer = setTimeout(function() {
            showText.value = "Weclome!";
        }, 3000)
    }

    function success(text: string) {
        setText(`[报告] ${text}`);
    }

    return { showText ,setText , success}
})
