<template>
  <v-tooltip :text="props.tooltipText" location="start">
    <template #activator="{ props }">
      <v-btn elevation="0" class="float-right" icon="mdi-pencil" size="small" v-bind="props" @click="dialogOpen = true"></v-btn>
    </template>
  </v-tooltip>

  <v-dialog v-model="dialogOpen" width="auto" persistent>
    <v-card>
      <template #title>
        <span v-if="props.dialogTitle" > {{ props.dialogTitle }} </span>
        <span v-else> {{ props.tooltipText }} </span>
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
    tooltipText: string,
    dialogTitle?: string
    onCancel?: () => Promise<string>,
    onSave?: () => Promise<string>,
  }>()

  const dialogOpen = ref(false)
  const errorMessage = ref("")

  function innerOnCancel() {
    if (props.onCancel) {
      props.onCancel().
        then(() => dialogOpen.value = false).
        catch((d) => errorMessage.value = d)
    } else {
      dialogOpen.value = false
    }
  }

  function innerOnSave() {
    if (props.onSave) {
      props.onSave().
        then(() => dialogOpen.value = false).
        catch((d) => errorMessage.value = d)
    } else {
      dialogOpen.value = false
    }
  }
</script>
