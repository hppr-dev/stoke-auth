<template>
  <v-data-table
    class="h-100"
    sticky
    :headerProps="headerProps"
    :headers="props.deleteClick? [ ...props.headers, { key: 'row-delete'} ] : props.headers "
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
    <template #item.row-icon="{ item }">
      <slot name="row-icon" :item="item"></slot>
    </template>
    <template v-if="!props.showFooter" #bottom></template>
    <template v-if="props.deleteClick" #item.row-delete="{ item }">
      <DeleteActivator
        :titleIcon="icons.USER"
        :deleteIcon="props.deleteIcon? props.deleteIcon: icons.DELETE"
        :onDelete="async () => await innerOnDelete(item)"
        :toDelete="item[props.deleteItemKey]"
      />
    </template>
  </v-data-table>
  <slot></slot>
</template>

<script setup lang="ts">
  import { ref, defineProps } from "vue"
  import icons from "../util/icons"
import DeleteActivator from "./DeleteActivator.vue";

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
    deleteIcon?: string,
    deleteClick?: Function,
    deleteItemKey?: string,
  }>()

  const search = ref("")

  const headerProps = {
    class : "bg-blue-grey",
    height: "3em",
  }

  async function innerOnDelete(item) {
    await props.deleteClick(item)
  }

</script>
