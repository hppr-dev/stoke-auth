// Utilities
import { defineStore } from 'pinia'
import { appActions } from './app_actions'
import { User, UserWithCreds, Group, Claim } from '../util/entityTypes'

interface PasswordForm {
  username : string
  oldPassword : string
  newPassword : string
  force : boolean
}

export const useAppStore = defineStore('app', {
  state: () => ({
    api_url: import.meta.env.DEV? import.meta.env.VITE_API_URL : "",
    username: "",
    token: "",
    refreshToken: "",
    refreshTimeout: 0,
    currentUser : {} as User,
    currentGroup: {} as Group,
    currentClaim: {} as Claim,

    scratchUser : {} as User | UserWithCreds,
    scratchGroup: {} as Group,
    scratchClaim: {} as Claim,

    currentGroups: [] as Group[],
    currentClaims: [] as Claim[],

    scratchGroups: [] as Group[],
    scratchClaims: [] as Claim[],

    allUsers:[] as User[],
    allGroups:[] as Group[],
    allClaims:[] as Claim[],

    passwordForm: {} as PasswordForm
  }),
  getters: {
    authenticated: function() {
      return this.token !== ""
    },
    tokenExpiration: function() {
      const b64 = this.token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')
      const claims = JSON.parse(decodeURIComponent(atob(b64)))
      return new Date(claims.exp * 1000)
    },
  },
  actions: appActions,
})
