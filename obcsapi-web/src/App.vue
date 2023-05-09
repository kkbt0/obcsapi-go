<script lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import HelloWorld from './components/HelloWorld.vue'
import { NConfigProvider, zhCN, dateZhCN, darkTheme, NMessageProvider } from 'naive-ui';
import { defineComponent } from 'vue';
import usemessageComponents from '@/components/api/usemessageComponents.vue';

export default defineComponent({
  components: {
    NConfigProvider,
    NMessageProvider, // N Provider
    usemessageComponents,
    HelloWorld, // View Components
  },
  setup() {
    return {
      darkTheme,
      zhCN,
      dateZhCN,
    }
  }
})

</script>

<template>
  <!-- Main Header -->
  <header>
    <div class="wrapper">
      <HelloWorld msg="Weclome!" />
      <nav>
        <RouterLink to="/">Home</RouterLink>
        <RouterLink to="/setting">Setting</RouterLink>
        <!-- <RouterLink to="/login">Login</RouterLink>
        <RouterLink to="/about">About</RouterLink> -->
      </nav>
    </div>
  </header>
  <!-- Main-->
  <n-config-provider :theme="darkTheme" :locale="zhCN" :date-locale="dateZhCN">
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
