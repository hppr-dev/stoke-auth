<template>
  <v-tooltip :text="props.tooltipText" location="start">
    <template #activator="{ props: activatorProps }">
      <v-btn :aria-label="props.tooltipText" elevation="0" class="float-right" :icon="icons.EDIT" size="small" v-bind="activatorProps" @click="dialogOpen = true"></v-btn>
    </template>
  </v-tooltip>

  <v-dialog v-model="dialogOpen" width="auto" @afterLeave="innerOnCancel" persistent>
    <v-alert
      class="mb-n16"
      type="error"
      :text="errorMsg"
      style="z-index : 100;"
      v-if="errorMsg"
      @click:close="errorMsg=''"
      closable
    ></v-alert>
    <v-form v-model="formValid" @submit.prevent="innerOnSave" validate-on="submit">
      <v-card>
        <template #title>
          <v-icon v-if="props.titleIcon" class="mr-2" :icon="props.titleIcon"></v-icon>
          <span v-if="props.dialogTitle" > {{ props.dialogTitle }} </span>
          <span v-else> {{ props.tooltipText }} </span>
          <v-divider></v-divider>
        </template>
        <template #text>
          <slot></slot>
        </template>
        <template #actions>
          <div class="d-flex justify-center w-100 pb-3">
            <v-btn color="error" @click="innerOnCancel">
              Cancel
            </v-btn>

            <v-btn color="success" type="submit">
              Save
            </v-btn>
          </div>
        </template>
      </v-card>
    </v-form>
  </v-dialog>
</template>

<script setup lang="ts">
  import { ref, defineProps } from 'vue'
  import icons from '../util/icons'

 const props = defineProps<{
    tooltipText: string,
    titleIcon?: string,
    dialogTitle?: string
    onCancel?: () => void,
    onSave?: () => Promise<string>,
  }>()

  const dialogOpen = ref(false)
  const formValid = ref(false)
  const errorMsg = ref("")

  async function hideDialog() {
    dialogOpen.value = false
    await new Promise(r => setTimeout(r, 250))
  }

  async function innerOnCancel() {
    await hideDialog()
    if(props.onCancel) {
      setTimeout(props.onCancel, 250)
    }
  }

  async function innerOnSave(event : Promise<SubmitEvent>) {
    await event
    if ( !formValid.value ) return

    if (props.onSave) {
      try {
        await props.onSave()
        await hideDialog()
      } catch (err) {
        console.error(err)
        if ( err.cause ) {
          errorMsg.value = err.cause
        } else {
          errorMsg.value = err + ''
        }
      }
    } else {
      await hideDialog()
    }
  }
</script>
