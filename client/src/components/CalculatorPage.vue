<template>
  <v-sheet
    height="500"
    class="fill-height d-flex justify-center flex-column mx-auto px-2 py-8"
  >
    <v-container max-width="500">
      <v-tabs
        v-model="tab"
        color="green"
        fixed-tabs
        class="rounded-t-lg mb-4"
      >
        <v-tab value="Calculator">Calculator</v-tab>
        <v-tab value="History">History</v-tab>
      </v-tabs>
      <v-tabs-window v-model="tab" class="px-2">
        <v-tabs-window-item class="pt-4" value="Calculator">
          <CalculatorForm/>
        </v-tabs-window-item>

        <v-tabs-window-item class="pt-2" value="History">
          <ExpressionsHistory/>
        </v-tabs-window-item>
      </v-tabs-window>
    </v-container>
  </v-sheet>
</template>

<script setup>
import CalculatorForm from "@/components/CalculatorForm.vue";
import {ref} from "vue";
import ExpressionsHistory from "@/components/ExpressionsHistory.vue";

const tab = ref(null)
</script>
