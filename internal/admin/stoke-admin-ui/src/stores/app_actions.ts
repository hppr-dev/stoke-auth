import { User, Claim, Group } from '../util/entityTypes'
import { parseMetricData } from '../util/prometheus'

export const appActions = {
  login: async function(username : string, password : string, callback : () => void) {
    try {
      const response = await fetch(`${this.api_url}/api/login`, {
        method: "POST",
        headers: {
          "Content-Type" : "application/json",
        },
        body: JSON.stringify({
          username: username,
          password: password,
          required_claims: {
            srol: "spr",
          },
        }),
      })

      if ( !response.ok ){
        throw new Error(response.statusText)
      }

      const result = await response.json();
      if ( result.message ) {
        throw new Error(result.message);
      }

      // Artificial wait to make it seem as robust as it is
      await new Promise((r) => setTimeout(r, 500))

      callback()

      this.username = username
      this.token = result.token
      this.refreshToken = result.refresh
      sessionStorage.setItem("token", result.token)
      sessionStorage.setItem("refresh", result.refresh)
      sessionStorage.setItem("username", username)

      this.scheduleRefresh()

    } catch (err) {

      this.username = ""
      this.token = ""
      this.refreshToken = ""
      this.refreshTimeout = 0
      sessionStorage.setItem("token", "")
      sessionStorage.setItem("refresh", "")
      sessionStorage.setItem("username", "")

      throw err
    }
  },
  refreshSession: async function() {
    try {
      const response = await fetch(`${this.api_url}/api/refresh`, {
        method: "POST",
        headers: {
          "Content-Type" : "application/json",
          "Authorization" : `Token ${this.token}`,
        },
        body: JSON.stringify({
          refresh: this.refreshToken,
        }),
      })

      if ( !response.ok ){
        throw new Error(response.statusText)
      }

      const result = await response.json();
      if ( result.message ) {
        throw new Error(result.message);
      }

      this.token = result.token
      this.refreshToken = result.refresh
      sessionStorage.setItem("token", result.token)
      sessionStorage.setItem("refresh", result.refresh)

      this.scheduleRefresh()

    } catch (err) {
      console.error(err)

      this.username = ""
      this.token = ""
      this.refreshToken = ""
      this.refreshTimeout = 0
      sessionStorage.setItem("token", "")
      sessionStorage.setItem("refresh", "")
      sessionStorage.setItem("username", "")

      throw err
    }
  },
  scheduleRefresh: function() {
    this.refreshTimeout = window.setTimeout(this.refreshSession, this.tokenExpiration.getTime() - Date.now() - 10000)
  },
  logout: function() {
    this.username = ""
    this.token = ""
    this.refreshToken = ""
    sessionStorage.setItem("token", "")
    sessionStorage.setItem("refresh", "")
    sessionStorage.setItem("username", "")

    clearTimeout(this.refreshTimeout)
  },
  simpleGet: async function(endpoint : string, stateName : string, refresh? : boolean) {
    if ( !refresh ) {
      return
    }

    const response = await fetch(`${this.api_url}${endpoint}`, {
      method: "GET",
      headers: {
        "Authorization" : `Token ${this.token}`,
      }
    })

    if ( !response.ok ){
      throw new Error(response.statusText)
    }

    this[stateName] = await response.json()
  },
  fetchAllUsers: async function(refresh? : boolean) {
    await this.simpleGet("/api/admin/users", "allUsers", this.allUsers.length == 0 || refresh)
  },
  fetchAllGroups: async function(refresh? : boolean) {
    await this.simpleGet("/api/admin/claim-groups", "allGroups", this.allGroups.length == 0 || refresh)
  },
  fetchGroupsForUser: async function(userId: number) {
    await this.simpleGet(`/api/admin/users/${userId}/claim-groups`, "currentGroups", true)
  },
  fetchAllClaims: async function(refresh? : boolean) {
    await this.simpleGet("/api/admin/claims", "allClaims", this.allClaims.length == 0 || refresh)
  },
  fetchClaimsForGroup: async function(groupId: number) {
    await this.simpleGet(`/api/admin/claim-groups/${groupId}/claims`, "currentClaims", true)
  },
  simplePatch: async function(endpoint : string, stateToSend : string) {
    const value : User | Claim | Group = this[stateToSend]
    const response = await fetch(`${this.api_url}${endpoint}/${value.id}`, {
      method: "PATCH",
      headers: {
        "Content-Type"  : "application/json",
        "Authorization" : `Token ${this.token}`,
      },
      body : JSON.stringify(value),
    })

    if ( !response.ok ){
      throw new Error(response.statusText)
    }
  },
  saveScratchUser: function() {
    this.currentUser = { ...this.scratchUser }
    this.currentGroups = [ ...this.scratchGroups ]
    this.scratchUser.claim_groups = this.scratchGroups.map((g : Group) => g.id)
    return this.simplePatch("/api/admin/users", "scratchUser")
      .then(() => this.scratchUser = {})
  },
  saveScratchGroup: function() {
    this.currentGroup = { ...this.scratchGroup }
    this.currentClaims = [ ...this.scratchClaims ]
    this.scratchGroup.claims = this.scratchClaims.map((c : Claim) => c.id)
    return this.simplePatch("/api/admin/claim-groups", "scratchGroup")
      .then(() => this.scratchGroup = {})
  },
  saveScratchClaim: function() {
    this.currentClaim = this.scratchClaim
    return this.simplePatch("/api/admin/claims", "scratchClaim")
      .then(() => this.scratchClaim = {})
  },
  savePasswordForm: async function() {
    const response = await fetch(`${this.api_url}/api/admin_users`, {
      method: "PATCH",
      headers: {
        "Content-Type"  : "application/json",
        "Authorization" : `Token ${this.token}`,
      },
      body : JSON.stringify(this.passwordForm),
    })
    if ( !response.ok ){
      throw new Error(response.statusText)
    }
  },
  simplePost: async function(endpoint : string, stateToSend : string) {
    const response = await fetch(`${this.api_url}${endpoint}`, {
      method: "POST",
      headers: {
        "Content-Type"  : "application/json",
        "Authorization" : `Token ${this.token}`,
      },
      body : JSON.stringify(this[stateToSend]),
    })
    if ( !response.ok ){
      throw new Error(response.statusText)
    }
  },
  addScratchUser: function() {
    return this.simplePost("/api/admin_users", "scratchUser")
      .then( () => this.scratchUser = {} )
      .then(this.fetchAllUsers)
  },
  addScratchGroup: function() {

    return this.simplePost("/api/admin/claim-groups", "scratchGroup")
      .then( () => this.scratchGroup = {} )
      .then(this.fetchAllGroups)
  },
  addScratchClaim: function() {
    return this.simplePost("/api/admin/claims", "scratchClaim")
      .then( () => this.scratchClaim = {} )
      .then(this.fetchAllClaims)
  },
  resetScratchUser: function() {
    this.scratchUser = {}
    this.scratchGroups = []
  },
  resetScratchGroup: function() {
    this.scratchGroup = {}
    this.scratchClaims = []
  },
  resetScratchClaim: function() {
    this.scratchClaim = {}
  },
  resetCurrentUser: function() {
    this.currentUser = {}
    this.currentGroups = []
  },
  resetCurrentGroup: function() {
    this.currentGroup = {}
    this.currentClaims = []
  },
  resetCurrentClaim: function() {
    this.currentClaim = {}
  },
  resetSelections: function() {
    this.scratchUser = {}
    this.scratchGroup = {}
    this.scratchClaim = {}

    this.scratchGroups = []
    this.scratchClaims = []

    this.currentUser = {}
    this.currentGroup = {}
    this.currentClaim = {}

    this.currentGroups = []
    this.currentClaims = []
  },
  fetchMetricData: async function() {
    const response = await fetch(`${this.api_url}/metrics`, {
        method: "GET",
        headers: {
          "Content-Type" : "text/plain; version=0.0.4",
          "Authorization" : `Token ${this.token}`,
        },
      }
    )
    const result = await response.text();
    this.metricData = parseMetricData(result);
  },
  metricRefresh: async function() {
    await this.fetchMetricData()
    this.metricTimeoutID = window.setTimeout(this.metricRefresh, this.metricRefreshTime)
    // TODO save data for charting
  },
  setMetricRefresh: function(millis : number) {
    if( millis < this.metricRefreshTime || this.metricTimeoutID === 0 ) {
      if( this.metricTimeoutID !== 0 ) {
        window.clearTimeout(this.metricTimeoutID)
      }
      this.metricTimeoutID = setTimeout(this.metricRefresh, millis)
    }
    this.metricRefreshTime = millis
  },
  clearMetricTimeout: function() {
    if( this.metricTimeoutID !== 0 ) {
      window.clearTimeout(this.metricTimeoutID)
    }
  }
}
