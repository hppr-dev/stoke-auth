import { User, Claim, Group, GroupLink } from '../util/entityTypes'

export const managementActions = {
  simpleGet: async function(endpoint : string, stateName : string, page: number, refresh? : boolean) {
    if ( !refresh && page == 1) {
      return
    }

    const response = await fetch(`${this.apiBase}${endpoint}?itemsPerPage=${this.pageLoadSize}&page=${page}`, {
      method: "GET",
      headers: {
        "Authorization" : `Bearer ${this.token}`,
      }
    })

    if ( !response.ok ){
      let cause = await response.json()
      throw new Error(response.statusText, { cause : cause.error_message })
    }

    const jres = await response.json()

    if ( page > 1 ) {
      this[stateName] = [...this[stateName], ...jres]
    } else {
      this[stateName] = jres
    }
  },
  fetchEntityTotals: async function() {
    await this.simpleGet("/admin/totals", "entityTotals", 1, true)
  },
  fetchAllUsers: async function(refresh? : boolean, page: number = 1) {
    await this.simpleGet("/admin/users", "allUsers", page, this.allUsers.length == 0 || refresh)
    await this.fetchEntityTotals()
  },
  fetchAllGroups: async function(refresh? : boolean, page: number = 1) {
    await this.simpleGet("/admin/claim-groups", "allGroups", page, this.allGroups.length == 0 || refresh)
    await this.fetchEntityTotals()
  },
  fetchAllClaims: async function(refresh? : boolean, page: number = 1) {
    await this.simpleGet("/admin/claims", "allClaims", page, this.allClaims.length == 0 || refresh)
    await this.fetchEntityTotals()
  },
  fetchGroupsForUser: async function(userId: number) {
    await this.simpleGet(`/admin/users/${userId}/claim-groups`, "currentGroups", 1, true)
  },
  fetchClaimsForGroup: async function(groupId: number) {
    await this.simpleGet(`/admin/claim-groups/${groupId}/claims`, "currentClaims", 1, true)
  },
  fetchLinksForGroup: async function(groupId: number) {
    await this.simpleGet(`/admin/claim-groups/${groupId}/group-links`, "currentLinks", 1, true)
  },
  simplePatch: async function(endpoint : string, stateToSend : string) {
    const value : User | Claim | Group = this[stateToSend]
    const response = await fetch(`${this.apiBase}${endpoint}/${value.id}`, {
      method: "PATCH",
      headers: {
        "Content-Type"  : "application/json",
        "Authorization" : `Bearer ${this.token}`,
      },
      body : JSON.stringify(value),
    })

    if ( !response.ok ){
      let cause = await response.json()
      throw new Error(response.statusText, { cause : cause.error_message })
    }
  },
  saveScratchUser: function() {
    this.scratchUser.claim_groups = this.scratchGroups.map((g : Group) => g.id)
    return this.simplePatch("/admin/users", "scratchUser")
      .then(() => {
        this.currentUser = { ...this.scratchUser }
        this.currentGroups = [ ...this.scratchGroups ]
        this.scratchUser = {}
      })
  },
  saveScratchGroup: function() {
    this.scratchGroup.claims = this.scratchClaims.map((c : Claim) => c.id)
    return this.simplePatch("/admin/claim-groups", "scratchGroup")
      .then(() => {
        this.currentGroup = { ...this.scratchGroup }
        this.currentClaims = [ ...this.scratchClaims ]
        this.scratchGroup = {}
      })
  },
  saveScratchClaim: function() {
    return this.simplePatch("/admin/claims", "scratchClaim")
      .then(() => {
        this.currentClaim = this.scratchClaim
        this.scratchClaim = {}
      })
  },
  savePasswordForm: async function() {
    const response = await fetch(`${this.apiBase}/admin/localuser`, {
      method: "PATCH",
      headers: {
        "Content-Type"  : "application/json",
        "Authorization" : `Bearer ${this.token}`,
      },
      body : JSON.stringify(this.passwordForm),
    })
    if ( !response.ok ){
      let cause = await response.json()
      throw new Error(response.statusText, { cause : cause.error_message })
    }
  },
  simplePost: async function(endpoint : string, stateToSend : string) {
    const response = await fetch(`${this.apiBase}${endpoint}`, {
      method: "POST",
      headers: {
        "Content-Type"  : "application/json",
        "Authorization" : `Bearer ${this.token}`,
      },
      body : JSON.stringify(this[stateToSend]),
    })
    if ( !response.ok ){
      let cause = await response.json()
      throw new Error(response.statusText, { cause : cause.error_message })
    }
  },
  addScratchUser: function() {
    return this.simplePost("/admin/localuser", "scratchUser")
      .then( () => {
        this.scratchUser = {}
        this.fetchAllUsers(true)
      })
  },
  addScratchGroup: function() {
    this.scratchGroup.claims = this.scratchClaims.map((c : Claim) => c.id)
    return this.simplePost("/admin/claim-groups", "scratchGroup")
      .then( () => {
        this.scratchGroup = {}
        this.fetchAllGroups(true)
      })
  },
  addScratchClaim: function() {
    return this.simplePost("/admin/claims", "scratchClaim")
      .then( () => {
        this.scratchClaim = {}
        this.fetchAllClaims(true)
      })
  },
  addScratchLink: function() {
    return this.simplePost("/admin/group-links", "scratchLink")
      .then( () => {
        this.scratchLink = {}
        this.fetchLinksForGroup(this.currentGroup.id)
      })
  },
  simpleDelete: async function(endpoint : string, value : User | Claim | Group | GroupLink) {
    const response = await fetch(`${this.apiBase}${endpoint}/${value.id}`, {
      method: "DELETE",
      headers: {
        "Content-Type"  : "application/json",
        "Authorization" : `Bearer ${this.token}`,
      },
      body : JSON.stringify(value),
    })

    if ( !response.ok ){
      let cause = await response.json()
      throw new Error(response.statusText, { cause : cause.error_message })
    }
  },
  deleteUser: function() {
    return this.simpleDelete("/admin/users", this.currentUser)
      .then(() => this.currentUser = {})
      .finally(() => this.fetchAllUsers(true))
  },
  deleteGroup: function() {
    return this.simpleDelete("/admin/claim-groups", this.currentGroup)
      .then(() => this.currentGroup = {})
      .finally(() => this.fetchAllGroups(true))
  },
  deleteClaim: function() {
    return this.simpleDelete("/admin/claims", this.currentClaim)
      .then(() => this.currentClaim = {})
      .finally(() => this.fetchAllClaims(true))
  },
  deleteLink: function(link: GroupLink) {
    return this.simpleDelete("/admin/group-links", link)
      .then(() => this.fetchLinksForGroup(this.currentGroup.id))
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
