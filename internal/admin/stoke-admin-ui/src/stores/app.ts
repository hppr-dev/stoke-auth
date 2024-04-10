// Utilities
import { defineStore } from 'pinia'
import { User, UserWithCreds, Group, Claim } from '../util/entityTypes'
import { MetricDataMap } from '../util/prometheus'
import { appActions } from './app_actions'
import { managementActions } from './management_actions'
import { metricActions, ChartData } from './metric_actions'

interface PasswordForm {
  username : string
  oldPassword : string
  newPassword : string
  force : boolean
}

export interface ChartDatasets {
  [k : string] : ChartData
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

    passwordForm: {} as PasswordForm,

    metricData: {} as MetricDataMap,
    metricsPaused: true,
    metricRefreshTime: 30000,
    metricTimeoutID: 0,
    maxPoints: 100,

    chartDatam: {} as ChartDatasets,
    logText: "",
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
    trackedMetrics: function() {
      return Object.keys(this.chartDatam)
    },
  },
  actions: {
    ...appActions,
    ...managementActions,
    ...metricActions,
  },
})
