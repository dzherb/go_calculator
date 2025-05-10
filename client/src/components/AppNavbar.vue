<template>
  <v-app-bar density="compact" elevation="0">
    <v-spacer/>
    <span class="mr-4 select-none">{{ currentUser?.username }}</span>
    <v-btn @click="logout()" title="Logout" icon="mdi-logout" variant="text" density="compact" class="mr-4"></v-btn>
  </v-app-bar>

</template>

<script setup>
import {useAuthentication} from "@/composables.js";

const {currentUser, logout} = useAuthentication()
</script>
