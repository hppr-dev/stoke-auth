<template>
  <v-btn class="ml-2 h-75" variant="tonal" color="error" stacked density="compact" @click="dialogOpen = true">
    <span>Change</span>
    <span>Password</span>
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
    <v-form v-model="formValid" @submit.prevent="onSave" validate-on="submit" ref="formRef">
      <v-card>
        <template #title>
          <v-icon class="mr-2" :icon="icons.USER"></v-icon>
          <span> Change Password </span>
          <v-divider></v-divider>
        </template>

        <template #text>
          <v-sheet class="px-4 py-4 mb-n5" width="40vw">
            <v-row>
              <v-col>
                <v-text-field
                  variant="solo-filled"
                  label="Username"
                  v-model="username"
                  :disabled="true"
                ></v-text-field>
              </v-col>
            </v-row>

            <v-row>
              <v-col cols="10">
                <v-text-field
                  variant="solo-filled"
                  label="Old Password"
                  no-resize
                  v-model="oldPassword"
                  :rules="[requireIfNotForce]"
                  :disabled="force"
                  :type="showOldPass? 'text' : 'password'"
                  :append-inner-icon="showOldPass ? icons.SHOW : icons.HIDE"
                  @click:append-inner="showOldPass = !showOldPass"
                  @update:modelValue="updatePasswordForm"
                ></v-text-field>
              </v-col>
              <v-col cols="2">
                <v-checkbox
                  v-model="force"
                  label="Force"
                  @click="onUpdateForce"
                ></v-checkbox>
              </v-col>
            </v-row>

            <v-row>
              <v-col>
                <v-text-field
                  type="password"
                  variant="solo-filled"
                  label="New Password"
                  no-resize
                  v-model="newPassword"
                  :rules="[require('New password'), matchPasswords]"
                  @update:modelValue="updatePasswordForm"
                ></v-text-field>
              </v-col>
              <v-col>
                <v-text-field
                  type="password"
                  variant="solo-filled"
                  label="Repeat New Password"
                  no-resize
                  v-model="newPasswordRepeat"
                  :rules="[require('Repeat new password'), matchPasswords]"
                  @update:modelValue="updatePasswordForm"
                ></v-text-field>
              </v-col>
            </v-row>

          </v-sheet>
        </template>

        <template #actions>
          <div class="d-flex justify-center w-100 pb-3">
            <v-btn color="error" @click="onCancel">
              Cancel
            </v-btn>

            <v-btn color="success" type="submit">
              Change Password
            </v-btn>
          </div>
        </template>
      </v-card>
    </v-form>
  </v-dialog>

</template>

<script setup lang="ts">
  import { ref, onMounted } from "vue"
  import { useAppStore } from "../../stores/app"
  import { require } from "../../util/rules"
  import icons from "../../util/icons"

  const store = useAppStore()

  const dialogOpen = ref(false)
  const formValid  = ref(false)
  const formRef    = ref(null)
  const showOldPass = ref(false)
  const errorMsg = ref("")

  const username = ref(store.currentUser.username)
  const oldPassword = ref("")
  const newPassword = ref("")
  const newPasswordRepeat = ref("")
  const force = ref(false)

  async function onSave(event : Promise<SubmitEvent>) {
    await event
    if ( !formValid.value ) {
      return
    }
    try {
      await store.savePasswordForm()
      dialogOpen.value = false
    } catch (err) {
      console.error(err)
      formRef.value.reset()
      errorMsg.value = "Could not update password. Re-enter passwords and try again."
    }
  }

  function onCancel() {
    dialogOpen.value = false
    store.$patch({
      passwordForm: {},
    })
  }

  function matchPasswords() {
    return newPassword.value === newPasswordRepeat.value || "New passwords do not match."
  }

  function requireIfNotForce() {
    return force.value || !!oldPassword.value || "Old password is required."
  }

  function onUpdateForce() {
    formRef.value.resetValidation()
    updatePasswordForm()
  }

  function updatePasswordForm() {
    store.$patch({
      passwordForm: {
        username: username.value,
        oldPassword: oldPassword.value,
        newPassword: newPassword.value,
        force: force.value,
      },
    })
  }

  onMounted(() => {
    store.$patch({
      passwordForm: {},
    })
  })
</script>
