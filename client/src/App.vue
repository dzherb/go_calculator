<template>
  <v-app>
    <AppNavbar/>
    <v-main>
      <CalculatorPage v-if="isAuthenticated"/>
      <AuthPage v-else/>
    </v-main>
    <AppFooter/>
  </v-app>
</template>

<script setup>
import AppFooter from "@/components/AppFooter.vue";
import {useAppTheme, useAuthentication} from "@/composables.js";
import AppNavbar from "@/components/AppNavbar.vue";
import {onBeforeMount} from "vue";
import CalculatorPage from "@/components/CalculatorPage.vue";
import AuthPage from "@/components/AuthPage.vue";

const {isAuthenticated, currentUser, fetchCurrentUser} = useAuthentication()

onBeforeMount(async () => {
  useAppTheme()

  if (isAuthenticated.value && currentUser.value === null) {
    await fetchCurrentUser()
  }
})
</script>
