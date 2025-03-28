export const appActions = {
  login: async function(username : string, password : string, provider : string, callback : () => void) {
    try {
      const response = await fetch(`${this.api_url}/api/login`, {
        method: "POST",
        headers: {
          "Content-Type" : "application/json",
        },
        body: JSON.stringify({
          username: username,
          password: password,
          provider: provider,
          // Match any stk claim
          required_claims: [ { "stk": "" } ],
          // Only return stk claims
          filter_claims: [ "stk" ],
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
      // and wait for any updates that need to happen
      await new Promise((r) => setTimeout(r, 500))

      callback()

      this.username = result.username
      this.token = result.token
      this.refreshToken = result.refresh
      sessionStorage.setItem("token", result.token)
      sessionStorage.setItem("refresh", result.refresh)
      sessionStorage.setItem("username", result.username)

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
          "Authorization" : `Bearer ${this.token}`,
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

      this.logout()
      throw err
    }
  },
  fetchCapabilites: async function() {
    const response = await fetch(`${this.api_url}/api/capabilities`, {
      method: "GET",
      headers: {
        "Content-Type" : "application/json",
        "Authorization" : `Bearer ${this.token}`,
      }
    })

    const result = await response.json();
    if ( result.capabilities ) {
      this.capabilites = result.capabilities
    }
  },
  fetchAvailableProviders: async function() {
    const response = await fetch(`${this.api_url}/api/available_providers`, {
      method: "GET",
      headers: {
        "Content-Type" : "application/json",
      }
    })

    const result = await response.json();
    if ( result.length > 0) {
      this.availableProviders = result
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

    window.location.replace("/admin/")
    clearTimeout(this.refreshTimeout)
  },
}
