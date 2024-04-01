
import { User, Claim, Group } from './entityTypes'

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
        }),
      })

      const result = await resp.json();
      if ( result.message ) {
        throw new Error(result.message);
      }

      this.username = username
      this.token = result.token
      sessionStorage.setItem("token", result.token)
      sessionStorage.setItem("username", username)

    } catch (err) {
      throw err
    }
  },
  logout: function() {
    this.username = ""
    this.token = ""
    sessionStorage.setItem("token", "")
    sessionStorage.setItem("username", "")
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
    return this.simplePost("/api/admin/users", "scratchUser")
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
}
