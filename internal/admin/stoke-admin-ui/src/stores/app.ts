// Utilities
import { defineStore } from 'pinia'
import { appActions } from './app_actions'
import { User, Group, Claim } from './entityTypes'

export const useAppStore = defineStore('app', {
  state: () => ({
    api_url: import.meta.env.DEV? import.meta.env.VITE_API_URL : "",
    username: "",
    token: "",
    currentUser : {} as User,
    currentGroup: {} as Group,
    currentClaim: {} as Claim,

    scratchUser : {} as User,
    scratchGroup: {} as Group,
    scratchClaim: {} as Claim,

    currentGroups: [] as Group[],
    currentClaims: [] as Claim[],

    scratchGroups: [] as Group[],
    scratchClaims: [] as Claim[],

    allUsers:[] as User[],
    allGroups:[] as Group[],
    allClaims:[] as Claim[],
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
