import { User, Claim, Group } from '../util/entityTypes'

export const managementActions = {
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
}
