<template>
  <v-btn class="mx-auto" color="success" prepend-icon="mdi-plus" @click="dialogOpen = true">
    {{ props.buttonText }}
  </v-btn>

  <v-dialog v-model="dialogOpen" width="auto" @afterLeave="onCancel" persistent>
    <v-card>
      <template #title>
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

          <v-btn color="success" @click="innerOnSave">
            Save
          </v-btn>
        </div>
      </template>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
  import { ref, defineProps } from 'vue'

 const props = defineProps<{
    buttonText: string,
    dialogTitle?: string
    onCancel?: () => Promise<string>,
    onSave?: () => Promise<string>,
  }>()

  const dialogOpen = ref(false)
  const errorMessage = ref("")

  function innerOnCancel() {
    dialogOpen.value = false
  }

  function innerOnSave() {
    if (props.onSave) {
      props.onSave().
        then(() => dialogOpen.value = false).
        catch((d) => { console.log(d) ; errorMessage.value = d })
    } else {
      dialogOpen.value = false
    }
  }
</script>
