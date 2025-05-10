<template>
  <v-container class="fill-height">
    <v-sheet
      width="300"
      height="500"
      class="align-centerfill-height mx-auto px-6 py-8"
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
        <p class="d-flex justify-space-between mt-2">Error: <span v-if="!error">-</span></p>
        <p class="text-red-darken-1 mt-2">{{ error }}</p>
      </v-form>
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
      return 'red'
  }
})
</script>

