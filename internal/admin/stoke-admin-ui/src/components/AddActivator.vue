<template>
  <v-btn class="mx-auto" color="success" :prepend-icon="icons.ADD" @click="dialogOpen = true">
    {{ props.buttonText }}
  </v-btn>

  <v-dialog v-model="dialogOpen" width="auto" @afterLeave="onCancel" persistent>
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
          <div class="d-flex justify-center w-100">
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
  const errorMessage = ref("")

  function innerOnCancel() {
    dialogOpen.value = false
  }

  async function innerOnSave(event : Promise<SubmitEvent>) {
    await event
    if ( !formValid.value ) return

    if (props.onSave) {
      try {
        await props.onSave()
        dialogOpen.value = false
      } catch (err) {
        console.error(err)
        errorMessage.value = err
      }
    } else {
      dialogOpen.value = false
    }
  }
</script>