<template>
  <v-data-table
    class="h-100"
    sticky
    :headerProps="headerProps"
    :headers="props.deleteClick? [ ...props.headers, { key: 'row-delete'} ] : props.headers "
    :items="props.items"
    :itemsPerPage="props.perPage"
    :search="search"
    :rowProps="rowProps"
    :loading="loading"
    v-model:page="page"
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


    <template #item.row-icon="{ item }">
      <slot name="row-icon" :item="item"></slot>
    </template>

    <template v-if="props.deleteClick" #item.row-delete="row">
      <DeleteActivator
        v-if="props.rowProps && props.rowProps(row)"
        :onDelete="props.deleteClick"
        :titleIcon="icons.USER"
        :toDelete="row.item[props.deleteItemKey]"
      />
    </template>

    <template #bottom>
      <div v-if="props.showFooter">
        <div class="text-center">
          <slot name="footer-prepend"></slot>
          <v-pagination
            v-model="page"
            :length="pageCount()"
            @next="innerOnNext"
          >
          </v-pagination>
        </div>
      </div>
      <div v-else></div>
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
    items: Array<object>,
    rowClick: Function,
    totalItems: number,
    perPage: number,
    onNext?: Function,
    searchIcon?: string,
    showSearch?: boolean,
    showFooter?: boolean,
    rowProps?: Function,
    deleteIcon?: string,
    deleteClick?: Promise<string>,
    deleteItemKey?: string,
  }>()

  const search = ref("")
  const page = ref(1)
  const loading = ref(false)

  const headerProps = {
    class : "bg-blue-grey",
    height: "3em",
  }

  async function innerOnNext(inPage: number) {
    loading.value = true
    if ( props.onNext ) {
      await props.onNext(inPage)
    }
    loading.value = false
  }

  async function innerOnDelete(item) {
    if ( props.deleteClick ) {
      await props.deleteClick(item)
    }
  }

  function pageCount() {
    return Math.ceil(props.totalItems / props.perPage)
  }

</script>
