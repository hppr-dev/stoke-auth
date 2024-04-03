import { User, Claim, Group } from '../util/entityTypes'

export const appActions = {
  login: async function(username : string, password : string) {
    try {
      const resp = await fetch(`${this.api_url}/api/login`, {
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

      const result = await resp.json();
      if ( result.message ) {
        throw new Error(result.message);
      }

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
      const resp = await fetch(`${this.api_url}/api/refresh`, {
        method: "POST",
        headers: {
          "Content-Type" : "application/json",
          "Authorization" : `Token ${this.token}`,
        },
        body: JSON.stringify({
          refresh: this.refreshToken,
        }),
      })

      const result = await resp.json();
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
    try {
      const resp = await fetch(`${this.api_url}${endpoint}`, {
        method: "GET",
        headers: {
          "Authorization" : `Token ${this.token}`,
        }
      })

      this[stateName] = await resp.json()
    } catch (err) {
      throw err
    }
  },
  fetchAllUsers: async function(refresh? : boolean) {
    try {
      await this.simpleGet("/api/admin/users", "allUsers", this.allUsers.length == 0 || refresh)
    } catch (err) {
      throw err
    }
  },
  fetchAllGroups: async function(refresh? : boolean) {
    try {
      await this.simpleGet("/api/admin/claim-groups", "allGroups", this.allGroups.length == 0 || refresh)
    } catch (err) {
      throw err
    }
  },
  fetchGroupsForUser: async function(userId: number) {
    try {
      await this.simpleGet(`/api/admin/users/${userId}/claim-groups`, "currentGroups", true)
    } catch (err) {
      throw err
    }
  },
  fetchAllClaims: async function(refresh? : boolean) {
    try {
      await this.simpleGet("/api/admin/claims", "allClaims", this.allClaims.length == 0 || refresh)
    } catch (err) {
      throw err
    }
  },
  fetchClaimsForGroup: async function(groupId: number) {
    try {
      await this.simpleGet(`/api/admin/claim-groups/${groupId}/claims`, "currentClaims", true)
    } catch (err) {
      throw err
    }
  },
  simplePatch: async function(endpoint : string, stateToSend : string) {
    try {
      const value : User | Claim | Group = this[stateToSend]
      await fetch(`${this.api_url}${endpoint}/${value.id}`, {
        method: "PATCH",
        headers: {
          "Content-Type"  : "application/json",
          "Authorization" : `Token ${this.token}`,
        },
        body : JSON.stringify(value),
      })
    } catch (err) {
      throw err
    }
  },
  saveScratchUser: function() {
    this.currentUser = { ...this.scratchUser }
    this.currentGroups = [ ...this.scratchGroups ]
    this.scratchUser.claim_groups = this.scratchGroups.map((g) => g.id)
    return this.simplePatch("/api/admin/users", "scratchUser")
      .then(() => this.scratchUser = {})
  },
  saveScratchGroup: function() {
    this.currentGroup = { ...this.scratchGroup }
    this.currentClaims = [ ...this.scratchClaims ]
    this.scratchGroup.claims = this.scratchClaims.map((c) => c.id)
    return this.simplePatch("/api/admin/claim-groups", "scratchGroup")
      .then(() => this.scratchGroup = {})
  },
  saveScratchClaim: function() {
    this.currentClaim = this.scratchClaim
    return this.simplePatch("/api/admin/claims", "scratchClaim")
      .then(() => this.scratchClaim = {})
  },
  simplePost: async function(endpoint : string, stateToSend : string) {
    try {
      await fetch(`${this.api_url}${endpoint}`, {
        method: "POST",
        headers: {
          "Content-Type"  : "application/json",
          "Authorization" : `Token ${this.token}`,
        },
        body : JSON.stringify(this[stateToSend]),
      })
    } catch (err) {
      throw err
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
  }
}
