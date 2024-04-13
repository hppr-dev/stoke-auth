<template>
  <v-btn class="mx-auto" color="success" :prepend-icon="icons.ADD" @click="dialogOpen = true">
    {{ props.buttonText }}
  </v-btn>

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
    <v-form v-model="formValid" @submit.prevent="innerOnSave" validate-on="submit">
      <v-card>
        <template #title>
          <v-icon v-if="props.titleIcon" class="mr-2" :icon="props.titleIcon"></v-icon>
          <span v-if="props.dialogTitle" > {{ props.dialogTitle }} </span>
          <span v-else> {{ props.buttonText }} </span>
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

            <v-btn type="submit" color="success">
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
    buttonText: string,
    titleIcon?: string,
    dialogTitle?: string
    onCancel?: () => Promise<string>,
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
      props.onCancel()
    }
  }

  async function innerOnSave(event : Promise<SubmitEvent>) {
    await event
    if ( !formValid.value ) return

    if (props.onSave) {
      try {
        await props.onSave()
        dialogOpen.value = false
        await hideDialog()
      } catch (err) {
        console.error(err)
        errorMsg.value = err + ''
      }
    } else {
      await hideDialog()
    }
  }
</script>
