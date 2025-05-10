<template>
  <v-chip v-if="status" density="compact" :color="statusColor">{{ status }}</v-chip>
</template>

<script setup>
import {computed} from "vue";
import {EXPRESSION_STATUS} from "@/api.js";

const props = defineProps(['status'])

const statusColor = computed(() => {
  switch (props.status) {
    case EXPRESSION_STATUS.NEW:
      return null
    case EXPRESSION_STATUS.PROCESSING:
      return 'primary'
    case EXPRESSION_STATUS.SUCCEED:
      return 'green'
    case EXPRESSION_STATUS.FAILED:
    case EXPRESSION_STATUS.ABORTED:
      return 'red-darken-4'
  }
})
</script>
