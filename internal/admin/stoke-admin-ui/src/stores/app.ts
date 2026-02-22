// Utilities
import { defineStore } from 'pinia'
import { User, UserWithCreds, Group, Claim, GroupLink } from '../util/entityTypes'
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

interface EntityTotals {
  users: number
  claims: number
  claim_groups: number
}

interface ProviderType {
  name: string
  provider_type: string
  type_spec: string
}

export interface ChartDatasets {
  [k : string] : ChartData
}

export const useAppStore = defineStore('app', {
  state: () => ({
    username: "",
    capabilities : [] as string[],

    token: "",
    refreshToken: "",
    refreshTimeout: 0,

    pageLoadSize: 200,

    entityTotals: {} as EntityTotals,
    currentUser : {} as User,
    currentGroup: {} as Group,
    currentClaim: {} as Claim,

    scratchUser : {} as User | UserWithCreds,
    scratchGroup: {} as Group,
    scratchClaim: {} as Claim,
    scratchLink:  {} as GroupLink,

    currentGroups: [] as Group[],
    currentClaims: [] as Claim[],
    currentLinks: [] as GroupLink[],

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

    availableProviders: [] as ProviderType[],
  }),
  getters: {
    /** Base URL for API requests: relative '../api' in production, or VITE_API_URL + '/api' in dev when set. */
    apiBase: function(): string {
      const env = import.meta.env
      if (env.DEV && env.VITE_API_URL) {
        const url = String(env.VITE_API_URL).replace(/\/$/, '')
        return url + '/api'
      }
      return '../api'
    },
    /** Base URL for metrics: relative '../metrics' in production, or VITE_API_URL + '/metrics' in dev when set. */
    metricsBase: function(): string {
      const env = import.meta.env
      if (env.DEV && env.VITE_API_URL) {
        const url = String(env.VITE_API_URL).replace(/\/$/, '')
        return url + '/metrics'
      }
      return '../metrics'
    },
    /** Server root for building URLs (e.g. OIDC): relative '..' in production, or VITE_API_URL in dev when set. */
    serverBase: function(): string {
      const env = import.meta.env
      if (env.DEV && env.VITE_API_URL) {
        return String(env.VITE_API_URL).replace(/\/$/, '')
      }
      return '..'
    },
    authenticated: function() {
      return this.token !== ""
    },
    tokenClaims: function() {
      const b64 = this.token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')
      return JSON.parse(decodeURIComponent(atob(b64)))
    },
    stokeClaim: function() : string[] {
      return this.tokenClaims.stk.split(",")
    },
    tokenExpiration: function() {
      return new Date(this.tokenClaims.exp * 1000)
    },
    superRead: function() : boolean {
      return this.stokeClaim.includes("s")
    },
    superUser: function() : boolean {
      return this.stokeClaim.includes("S")
    },
    monitoringAccess: function() : boolean {
      return this.capabilities.includes('monitoring') && (this.superUser() || this.superRead())
    },
    userAccess: function() : string {
      const stk = this.stokeClaim
      if ( this.superRead || stk.includes("u") ){
        return "read"
      } else if ( this.superUser || stk.includes("U") ) {
        return "write"
      }
      return ""
    },
    claimsAccess: function() : string {
      const stk = this.stokeClaim
      if ( this.superRead || stk.includes("c") ){
        return "read"
      } else if ( this.superUser || stk.includes("C") ) {
        return "write"
      }
      return ""
    },
    groupAccess: function() : string {
      const stk = this.stokeClaim
      if ( this.superRead || stk.includes("g") ){
        return "read"
      } else if ( this.superUser || stk.includes("G") ) {
        return "write"
      }
      return ""
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
