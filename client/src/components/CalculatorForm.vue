<template>
  <v-container class="fill-height">
    <v-sheet
      width="300"
      height="500"
      class="fill-height d-flex justify-center flex-column mx-auto px-6 py-8"
    >
      <v-form
        @submit.prevent
        :disabled="isLoading"
      >
        <h2 class="text-h6 pb-6">GO calculator</h2>
        <v-text-field
          label="Expression"
          v-model="expression"
          variant="outlined"
          density="compact"
          clearable
        />
        <v-btn
          :loading="isLoading"
          :disabled="isSendButtonDisabled"
          variant="flat"
          color="green"
          class="text-subtitle-2"
          @click="send()"
        >
          Send
        </v-btn>
        <v-divider class="my-6"/>
        <p class="d-flex justify-space-between align-center">
          Result:
          <span class="text-caption font-weight-bold" v-if="result">{{ result }}</span>
          <span v-else>-</span>
        </p>
        <p class="d-flex justify-space-between mt-2">
          Status:
          <v-chip v-if="status" density="compact" :color="statusColor">{{ status }}</v-chip>
          <span v-else>-</span>
        </p>
        <div class="error-container">
          <p class="mt-2" :class="error ? '' : 'd-flex justify-space-between'">
            Error: <span v-if="!error">-</span>
            <span v-if="error" class="text-red-darken-4 ml-1">{{ error }}</span>
          </p>
        </div>
      </v-form>
      <div class="h-25"></div>
    </v-sheet>
  </v-container>
</template>

<script setup>
import {computed, ref} from "vue";
import {useExpressionServerEvaluation} from "@/composables.js";
import {EXPRESSION_STATUS} from "@/api.js";

const expression = ref('')
const isSendButtonDisabled = computed(() => !expression.value)

const {result, status, error, isLoading, send} = useExpressionServerEvaluation(expression)

const statusColor = computed(() => {
  switch (status.value) {
    case EXPRESSION_STATUS.NEW:
      return null
    case EXPRESSION_STATUS.PROCESSING:
      return 'primary'
    case EXPRESSION_STATUS.SUCCEED:
      return 'green'
    case EXPRESSION_STATUS.FAILED || EXPRESSION_STATUS.ABORTED:
      return 'red-darken-4'
  }
})
</script>

