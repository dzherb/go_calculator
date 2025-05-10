<template>
  <v-app>
    <AppNavbar v-if="isAuthenticated"/>
    <v-main>
      <CalculatorForm v-if="isAuthenticated"/>
      <AuthForm v-else/>
    </v-main>
    <AppFooter />
  </v-app>
</template>

<script setup>
  import CalculatorForm from "@/components/CalculatorForm.vue";
  import AppFooter from "@/components/AppFooter.vue";
  import AuthForm from "@/components/AuthForm.vue";
  import {useAuthentication} from "@/composables.js";
  import AppNavbar from "@/components/AppNavbar.vue";
  import {onBeforeMount} from "vue";

  const {isAuthenticated, currentUser, fetchCurrentUser} = useAuthentication()

  onBeforeMount(async () => {
    if (isAuthenticated.value && currentUser.value === null) {
      await fetchCurrentUser()
    }
  })
</script>
