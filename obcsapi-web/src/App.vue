<script lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import HelloWorld from './components/HelloWorld.vue'
import { NConfigProvider, zhCN, dateZhCN, darkTheme, NMessageProvider } from 'naive-ui';
import { defineComponent, ref, onMounted } from 'vue';
import usemessageComponents from '@/components/api/usemessageComponents.vue';

export default defineComponent({
  components: {
    NConfigProvider,
    NMessageProvider, // N Provider
    usemessageComponents,
    HelloWorld, // View Components
  },

  setup() {
    const theme = ref();
    function switchTheme() {
      theme.value = theme.value == null ? darkTheme : null;
    }
    onMounted(() => {
      if (!window.matchMedia('(prefers-color-scheme: light)').matches) {
        switchTheme()
      }
      console.log("%cObcsapi"," text-shadow: 0 1px 0 #ccc,0 2px 0 #c9c9c9,0 3px 0 #bbb,0 4px 0 #b9b9b9,0 5px 0 #aaa,0 6px 1px rgba(0,0,0,.1),0 0 5px rgba(0,0,0,.1),0 1px 3px rgba(0,0,0,.3),0 3px 5px rgba(0,0,0,.2),0 5px 10px rgba(0,0,0,.25),0 10px 10px rgba(0,0,0,.2),0 20px 20px rgba(0,0,0,.15);font-size:5em")
      console.log("Doc  https://kkbt.gitee.io/obcsapi-go/#/md/go-version")
      console.log("v20230510-0853")
    
    })
    return {
      theme,
      darkTheme,
      zhCN,
      dateZhCN,
      switchTheme
    }
  }
})

</script>

<template>
  <!-- Main Header -->
  <header>
    <div class="wrapper">
      <HelloWorld msg="Weclome!" @click="" />
      <nav>
        <RouterLink to="/">Home</RouterLink>
        <RouterLink to="/setting">Setting</RouterLink>
        <!-- <a @click="switchTheme">🌞</a> -->
        <!-- <RouterLink to="/login">Login</RouterLink>
        <RouterLink to="/about">About</RouterLink> -->
      </nav>
    </div>
  </header>
  <!-- Main-->
  <n-config-provider :theme="theme" :locale="zhCN" :date-locale="dateZhCN">
    <!-- Info Components -->
    <n-message-provider>
      <usemessageComponents />
    </n-message-provider>
    <RouterView />
  </n-config-provider>
</template>

<style scoped>
header {
  line-height: 1.5;
  max-height: 100vh;
}

.logo {
  display: block;
  margin: 0 auto 2rem;
}

nav {
  width: 100%;
  font-size: 12px;
  text-align: center;
  margin-top: 3px;
}

nav a.router-link-exact-active {
  color: var(--color-text);
}

nav a.router-link-exact-active:hover {
  background-color: transparent;
}

nav a {
  display: inline-block;
  padding: 0 1rem;
  border-left: 1px solid var(--color-border);
}

nav a:first-of-type {
  border: 0;
}

@media (min-width: 1024px) {
  header {
    display: flex;
    place-items: center;
    padding-right: calc(var(--section-gap) / 2);
  }

  .logo {
    margin: 0 2rem 0 0;
  }

  header .wrapper {
    display: flex;
    place-items: flex-start;
    flex-wrap: wrap;
  }

  nav {
    text-align: left;
    margin-left: -1rem;
    font-size: 1rem;

    padding: 1rem 0;
    margin-top: 1rem;
  }
}
</style>
