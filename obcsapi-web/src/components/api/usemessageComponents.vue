<!-- Vue3 Message naive ui useMessage -->
<script lang="ts">
import { useMessage } from 'naive-ui'
import { defineComponent } from 'vue'

// content
export default defineComponent({
  setup () {
    window.$message = useMessage()
  }
})
</script>
<template>
  
</template>