<template>
  <v-card elevation="0">
    <v-list class="d-flex justify-space-between align-center px-2 pt-2 pb-6">
      <span class="text-subtitle-2 basis-full text-center">
        Expression
      </span>
      <span class="text-subtitle-2 basis-full text-center">
        Result
      </span>
      <span class="text-subtitle-2 basis-full text-center">
        Status
      </span>
    </v-list>
    <v-list height="255">
      <template v-for="expression in expressions">
        <div class="d-flex justify-space-between align-center px-2 py-2">
          <span class="text-subtitle-2 basis-full text-center">
            {{ expression.expression }}
          </span>
          <span class="text-subtitle-2 basis-full text-center">
            {{ expression.result || '-' }}
          </span>
          <span class="basis-full text-center">
            <StatusChip :status="expression.status"/>
          </span>
        </div>
        <v-divider></v-divider>
      </template>
    </v-list>
    <v-list class="d-flex justify-end">
      <v-btn
        :disabled="isLoading"
        :loading="isLoading"
        prepend-icon="mdi-reload"
        density="comfortable"
        variant="flat"
        color="green"
        class="text-subtitle-2 mt-4"
        @click="fetchExpressions"
      >
        Reload
      </v-btn>
    </v-list>
  </v-card>
</template>

<script setup>

import {useExpressionsHistory} from "@/composables.js";
import {onBeforeMount} from "vue";
import StatusChip from "@/components/StatusChip.vue";

const {expressions, isLoading, error, fetchExpressions} = useExpressionsHistory()

onBeforeMount(async () => await fetchExpressions())
</script>

<style scoped>
.basis-full {
  flex-basis: 100%;
}
</style>
