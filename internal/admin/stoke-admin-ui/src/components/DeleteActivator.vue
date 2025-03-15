<template>
  <v-tooltip :text="'Delete ' + compProps.toDelete" location="start">
    <template #activator="{ props }">
      <v-btn
        elevation="0"
        color="error"
        size="small"
        variant="text"
        :icon="compProps.deleteIcon? compProps.deleteIcon: icons.DELETE"
        v-bind="props"
        @click.stop="dialogOpen = true"
      ></v-btn>
    </template>
  </v-tooltip>

  <v-dialog v-model="dialogOpen" width="auto" @afterLeave="onCancel" persistent>
    <v-alert
      class="mb-n16"
      type="error"
      :text="errorMsg"
      style="z-index : 100;"
      v-if="errorMsg"
      @click:close="errorMsg=''"
      closable
    ></v-alert>
    <v-card>
      <template #title>
        <div class="pb-2 text-center">
          <v-icon v-if="compProps.titleIcon" class="mr-2" :icon="compProps.titleIcon"></v-icon>
          <span > Delete {{ compProps.toDelete }}? </span>
        </div>
        <v-divider></v-divider>
      </template>
      <template #text>
        <slot>
          <div class="py-4 text-center">
            <p>
              Are you sure you want to delete '{{compProps.toDelete}}'?
            </p>
            <p>
              This action can not be undone.
            </p>
          </div>
        </slot>
      </template>
      <template #actions>
        <div class="d-flex justify-end w-100 pb-3">
          <v-btn class="mx-2" color="primary" variant="elevated" @click="innerOnCancel">
            Cancel
          </v-btn>

          <v-btn class="mx-2" color="error" variant="elevated" :append-icon="compProps.deleteIcon? compProps.deleteIcon: icons.DELETE" @click="innerOnDelete">
            Delete
          </v-btn>
        </div>
      </template>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
  import { ref, defineProps } from 'vue'
  import icons from '../util/icons'

 const compProps = defineProps<{
    toDelete: string,
    titleIcon?: string,
    deleteIcon?: string,
    onCancel?: () => void,
    onDelete?: () => Promise<string>,
  }>()

  const dialogOpen = ref(false)
  const errorMsg = ref("")

  function innerOnCancel() {
    dialogOpen.value = false
    if(compProps.onCancel) {
      compProps.onCancel()
    }
  }

  async function innerOnDelete() {
    if (compProps.onDelete) {
      try {
        await compProps.onDelete()
        dialogOpen.value = false
      } catch (err) {
        console.error(err)
        if ( err.cause ) {
          errorMsg.value = err.cause
        } else {
          errorMsg.value = err + ''
        }
      }
    } else {
      dialogOpen.value = false
    }
  }
</script>
