<template>
  <v-data-table
    class="h-100"
    :headers="props.headers"
    :items="props.items"
    :search="search"
    @click:row="rowClick"
  >
    <template #top>
      <div v-if="props.showSearch">
        <v-text-field
          v-model="search"
          label="Search"
          prepend-inner-icon="mdi-magnify"
          variant="outlined"
          hide-details
          single-line
        ></v-text-field>
      </div>
    </template>
    <template #footer.prepend>
      <slot name="footer-prepend"></slot>
    </template>
    <template v-if="!props.showFooter" #bottom></template>
  </v-data-table>
  <slot></slot>
</template>

<script setup lang="ts">
  import { ref, defineProps } from "vue"

  interface Headers {
    key: string,
    title: string,
  }

  const props = defineProps<{
    showSearch?: boolean,
    showFooter?: boolean,
    headers : Headers,
    items: Array<Object>,
    rowClick?: Function,
  }>()

  const search = ref("")

</script>
