<template>
  <v-app-bar density="compact" elevation="0">
    <v-btn
      class="ml-4"
      color="green"
      title="Toggle theme"
      @click="toggleTheme"
      icon="mdi-brightness-4"
      density="compact">
    </v-btn>
    <template v-if="isAuthenticated">
      <v-spacer/>
      <span class="mr-4 select-none">{{ currentUser?.username }}</span>
      <v-btn
        @click="logout()"
        title="Logout"
        icon="mdi-logout"
        variant="text"
        density="compact"
        class="mr-4">
      </v-btn>
    </template>
  </v-app-bar>
</template>

<script setup>
import {useAppTheme, useAuthentication} from "@/composables.js";

const {isAuthenticated, currentUser, logout} = useAuthentication()

const {toggleTheme} = useAppTheme()
</script>
