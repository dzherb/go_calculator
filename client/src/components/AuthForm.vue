<template>
  <v-container class="fill-height">
    <v-sheet
      width="350"
      height="500"
      class="fill-height d-flex justify-center flex-column fill-height mx-auto px-6 py-8"
    >
      <v-form
        v-model="isFormValid"
        @submit.prevent
      >
        <v-tabs
          :disabled="isLoading"
          v-model="tab"
          color="green"
          fixed-tabs
          class="rounded-t-lg"
        >
          <v-tab value="Login">Login</v-tab>
          <v-tab value="Register">Register</v-tab>
        </v-tabs>
        <v-text-field
          :disabled="isLoading"
          :rules="[required, minLength(4)]"
          label="Username"
          v-model="username"
          variant="outlined"
          density="compact"
          clearable
          class="mt-8"
        />
        <v-text-field
          :disabled="isLoading"
          :rules="[required, minLength(8)]"
          label="Password"
          v-model="password"
          type="password"
          variant="outlined"
          density="compact"
          clearable
          class="mt-2"
        />
        <div class="d-flex justify-end">
          <v-btn
            :disabled="isLoading || !isFormValid"
            :loading="isLoading"
            variant="flat"
            color="green"
            class="text-subtitle-2 mt-2"
            @click="action()"
          >
            {{ tab }}
          </v-btn>
        </div>
        <p class="text-red-darken-4 mt-2 error-container">{{ loginError || registerError }}</p>
      </v-form>
      <div class="h-25"></div>
    </v-sheet>
  </v-container>
</template>

<script setup>
import {ref, toValue, watchEffect} from "vue";
import {useAuthentication} from "@/composables.js";
import {timeout} from "@/utils.js";

const isFormValid = ref(false)

const tab = ref(null)

const username = ref('')
const password = ref('')

const isLoading = ref(false)

const {loginError, registerError, login, register} = useAuthentication()

const required = (v) => {
  return !!v || 'Field is required'
}

const minLength = (length) => {
  return (v) => v?.length && v.length >= length || `Field must be at least ${length} characters`
}

const action = async () => {
  isLoading.value = true
  await timeout(3000)
  try {
    if (tab.value === 'Login') {
      await login({username, password})
      return
    }

    await register({username, password})
  } finally {
    isLoading.value = false
  }
}

watchEffect(() => {
  toValue(tab)
  toValue(username)
  toValue(password)
  loginError.value = null
  registerError.value = null
})

</script>
