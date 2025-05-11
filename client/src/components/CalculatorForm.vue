<template>
  <v-form
    @submit.prevent
    :disabled="isLoading"
  >
    <v-text-field
      :counter="MAX_LENGTH"
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
      Calculate
    </v-btn>
    <v-divider class="my-6"/>
    <p class="d-flex justify-space-between align-center">
      Result:
      <span class="text-caption font-weight-bold" v-if="result">{{ result }}</span>
      <span v-else>-</span>
    </p>
    <p class="d-flex justify-space-between mt-2">
      Status:
      <StatusChip v-if="status" :status/>
      <span v-else>-</span>
    </p>
    <div class="error-container">
      <p class="mt-2" :class="error ? '' : 'd-flex justify-space-between'">
        Error: <span v-if="!error">-</span>
        <span v-if="error" class="text-red-darken-4 ml-1">{{ error }}</span>
      </p>
    </div>
  </v-form>
  <div style="height: 100px"></div>
</template>

<script setup>
import {computed, ref, watch, watchEffect} from "vue";
import {useExpressionServerEvaluation, useExpressionsHistory} from "@/composables.js";
import StatusChip from "@/components/StatusChip.vue";

const MAX_LENGTH = 256

const expression = ref('')
const isSendButtonDisabled = computed(
  () => !expression.value || !expression.value.trim() || expression.value.length > MAX_LENGTH
)

const {result, status, error, isLoading, send} = useExpressionServerEvaluation(expression)
</script>
