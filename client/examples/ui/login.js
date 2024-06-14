const api_url = "http://localhost:8080"

class StokeTokenManager {
	async login(user, password, required_claims, on_refresh) {
		const response = await fetch(`${api_url}/api/login`, {
			method: "POST",
			headers: {
				"Content-Type" : "application/json",
			},
			body: JSON.stringify({
				username: user,
				password: password,
				required_claims: required_claims,
			}),
		})

		if ( !response.ok ){
			throw new Error(response.statusText)
		}

		const result = await response.json();
		if ( result.message ) {
			throw new Error(result.message);
		}

		this.on_refresh = on_refresh

		this.username = user
		this.token = result.token
		this.refreshToken = result.refresh
		this.claims = this.parseClaims()
		this.tokenExpiration = new Date(this.claims.exp * 1000)

		sessionStorage.setItem("username", this.username)
		sessionStorage.setItem("token", this.token)
		sessionStorage.setItem("refresh", this.refreshToken)

		this.on_refresh()

		this.scheduleRefresh()
	}

	async refreshSession() {
		const response = await fetch(`${api_url}/api/refresh`, {
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
		this.claims = this.parseClaims()
		this.tokenExpiration = new Date(this.claims.exp * 1000)

		sessionStorage.setItem("token", this.token)
		sessionStorage.setItem("refresh", this.refreshToken)

		this.scheduleRefresh()
	}

	async makeRequest(url, method, body_obj) {
		let req_data = {
			method: method,
			headers: {
				"Content-Type" : "application/json",
				"Authorization" : `Bearer ${this.token}`,
			},
		}
		if ( method !== "GET" ) {
			req_data.body = JSON.stringify(body_obj)
		}
		return await fetch(url, req_data)
	}

	scheduleRefresh() {
		const millis= this.tokenExpiration.getTime() - (new Date().getTime() + 1000);
		this.refreshTimeout = window.setTimeout(async () => {
			await this.refreshSession()
			this.on_refresh()
		}, millis) 
	}

	logout() {
		this.username = ""
		this.token = ""
		this.refreshToken = ""
		sessionStorage.setItem("username", null)
		sessionStorage.setItem("token", null)
		sessionStorage.setItem("refresh", null)

		window.clearTimeout(this.refreshTimeout)
	}

	parseClaims() {
		const b64 = this.token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')
		const claims = JSON.parse(decodeURIComponent(atob(b64)))
		return claims
	}
}
