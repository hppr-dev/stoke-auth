// Utilities
import { defineStore } from 'pinia'
import { appActions } from './app_actions'
import { User, Group, Claim } from './entityTypes'

export const useAppStore = defineStore('app', {
  state: () => ({
    api_url: import.meta.env.DEV? import.meta.env.VITE_API_URL : "",
    username: "",
    token: "",
    currentClaim: {} as Claim,
    currentClaims: [] as Claim[],
    allClaims:[] as Claim[],
    currentGroup: {} as Group,
    currentGroups:[] as Group[],
    allGroups:[] as Group[],
    currentUser : {} as User,
    allUsers:[] as User[],
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
