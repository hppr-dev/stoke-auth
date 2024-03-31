<template>
  <v-footer height="40" app>
    <v-switch class="mt-5" color="error" @click="updateTheme" v-model="themeSelect"></v-switch>
    <div
      class="text-caption text-disabled"
      style="position: absolute; right: 50px;"
    >
      &copy; 2016-{{ (new Date()).getFullYear() }}
      <a
        class="text-decoration-none on-surface"
        href="https://hppr.dev"
        rel="noopener noreferrer"
        target="_blank"
      >
      hppr.dev
      </a>
    </div>
  </v-footer>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { useTheme } from "vuetify"

const theme = useTheme()
const themeSelect = ref(false)

function updateTheme() {
  theme.global.name.value = themeSelect.value? "dark" : "light"
}
</script>
