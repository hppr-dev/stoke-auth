<template>
  <v-tooltip :text="props.tooltipText" location="start">
    <template #activator="{ props }">
      <v-btn elevation="0" class="float-right" icon="mdi-pencil" size="small" v-bind="props" @click="dialogOpen = true"></v-btn>
    </template>
  </v-tooltip>

  <v-dialog v-model="dialogOpen" width="auto" @afterLeave="onCancel" persistent>
    <v-form v-model="formValid" @submit.prevent="innerOnSave" validate-on="submit">
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

 const props = defineProps<{
    tooltipText: string,
    dialogTitle?: string
    onCancel?: () => void,
    onSave?: () => Promise<string>,
  }>()

  const dialogOpen = ref(false)
  const formValid = ref(false)
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
