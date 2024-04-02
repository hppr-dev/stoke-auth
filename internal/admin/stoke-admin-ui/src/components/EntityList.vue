<template>
  <v-data-table
    class="h-100"
    sticky
    :headerProps="headerProps"
    :headers="props.headers"
    :items="props.items"
    :search="search"
    :rowProps="rowProps"
    @click:row="rowClick"
  >
    <template #top>
      <div v-if="props.showSearch">
        <v-text-field
          v-model="search"
          label="Search"
          :prepend-inner-icon="props.searchIcon"
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
    headers : Headers,
    items: Array<Object>,
    rowClick: Function,
    searchIcon?: string,
    showSearch?: boolean,
    showFooter?: boolean,
    rowProps?: Function | Object,
  }>()

  const search = ref("")

  const headerProps = {
    class : "bg-blue-grey",
    height: "3em",
  }

</script>
