import { User, Claim, Group } from '../util/entityTypes'

export const managementActions = {
  simpleGet: async function(endpoint : string, stateName : string, page: number, refresh? : boolean) {
    if ( !refresh && page == 1) {
      return
    }

    const response = await fetch(`${this.api_url}${endpoint}?itemsPerPage=${this.pageLoadSize}&page=${page}`, {
      method: "GET",
      headers: {
        "Authorization" : `Bearer ${this.token}`,
      }
    })

    if ( !response.ok ){
      throw new Error(response.statusText)
    }

    const jres = await response.json()

    if ( page > 1 ) {
      this[stateName] = [...this[stateName], ...jres]
    } else {
      this[stateName] = jres
    }
  },
  fetchEntityTotals: async function() {
    await this.simpleGet("/api/admin/totals", "entityTotals", 1, true)
  },
  fetchAllUsers: async function(refresh? : boolean, page: number = 1) {
    await this.simpleGet("/api/admin/users", "allUsers", page, this.allUsers.length == 0 || refresh)
    await this.fetchEntityTotals()
  },
  fetchAllGroups: async function(refresh? : boolean, page: number = 1) {
    await this.simpleGet("/api/admin/claim-groups", "allGroups", page, this.allGroups.length == 0 || refresh)
    await this.fetchEntityTotals()
  },
  fetchAllClaims: async function(refresh? : boolean, page: number = 1) {
    await this.simpleGet(`/api/admin/claims`, "allClaims", page, this.allClaims.length == 0 || refresh)
    await this.fetchEntityTotals()
  },
  fetchGroupsForUser: async function(userId: number) {
    await this.simpleGet(`/api/admin/users/${userId}/claim-groups`, "currentGroups", 1, true)
  },
  fetchClaimsForGroup: async function(groupId: number) {
    await this.simpleGet(`/api/admin/claim-groups/${groupId}/claims`, "currentClaims", 1, true)
  },
  fetchLinksForGroup: async function(groupId: number) {
    await this.simpleGet(`/api/admin/claim-groups/${groupId}/group-links`, "currentLinks", 1, true)
  },
  simplePatch: async function(endpoint : string, stateToSend : string) {
    const value : User | Claim | Group = this[stateToSend]
    const response = await fetch(`${this.api_url}${endpoint}/${value.id}`, {
      method: "PATCH",
      headers: {
        "Content-Type"  : "application/json",
        "Authorization" : `Bearer ${this.token}`,
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
    const response = await fetch(`${this.api_url}/api/admin/localuser`, {
      method: "PATCH",
      headers: {
        "Content-Type"  : "application/json",
        "Authorization" : `Bearer ${this.token}`,
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
        "Authorization" : `Bearer ${this.token}`,
      },
      body : JSON.stringify(this[stateToSend]),
    })
    if ( !response.ok ){
      throw new Error(response.statusText)
    }
  },
  addScratchUser: function() {
    return this.simplePost("/api/admin/localuser", "scratchUser")
      .then( () => this.scratchUser = {} )
      .then(this.fetchAllUsers)
  },
  addScratchGroup: function() {
    this.scratchGroup.claims = this.scratchClaims.map((c : Claim) => c.id)
    return this.simplePost("/api/admin/claim-groups", "scratchGroup")
      .then( () => this.scratchGroup = {} )
      .then(this.fetchAllGroups)
  },
  addScratchClaim: function() {
    return this.simplePost("/api/admin/claims", "scratchClaim")
      .then( () => this.scratchClaim = {} )
      .then(this.fetchAllClaims)
  },
  addScratchLink: function() {
    return this.simplePost("/api/admin/group-links", "scratchLink")
      .then( () => {
        this.scratchLink = {}
        this.fetchLinksForGroup(this.currentGroup.id)
      })
  },
  simpleDelete: async function(endpoint : string, value : User | Claim | Group) {
    const response = await fetch(`${this.api_url}${endpoint}/${value.id}`, {
      method: "DELETE",
      headers: {
        "Content-Type"  : "application/json",
        "Authorization" : `Bearer ${this.token}`,
      },
      body : JSON.stringify(value),
    })

    if ( !response.ok ){
      throw new Error(response.statusText)
    }
  },
  deleteUser: function(user: User) {
    return this.simpleDelete("/api/admin/users", user)
      .then(() => this.currentUser = {})
      .then(this.fetchAllUsers)
  },
  deleteGroup: function(group: Group) {
    return this.simpleDelete("/api/admin/claim-groups", group)
      .then(() => this.currentGroup = {})
      .then(this.fetchAllGroups)
  },
  deleteClaim: function(claim: Claim) {
    return this.simpleDelete("/api/admin/claims", claim)
      .then(() => this.currentClaim = {})
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

    this.resetCurrentGroup()
  },
  resetCurrentGroup: function() {
    this.currentGroup = {}
    this.currentClaims = []

    this.resetCurrentClaim()
  },
  resetCurrentClaim: function() {
    this.currentClaim = {}
  },
  resetCurrentLink: function() {
    this.currentLink = {}
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
