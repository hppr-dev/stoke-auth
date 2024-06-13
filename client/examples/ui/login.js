const api_url = "http://localhost:8080"

class StokeTokenManager {
	async login(user, password, required_claims) {
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


		this.username = user
		this.token = result.token
		this.refreshToken = result.refresh
		this.tokenExpiration = this.computeExpiration()

		sessionStorage.setItem("token", this.token)
		sessionStorage.setItem("username", this.username)
		sessionStorage.setItem("tokenExpiration", this.tokenExpiration)

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
		this.tokenExpiration = this.computeExpiration()

		sessionStorage.setItem("token", this.token)
		sessionStorage.setItem("tokenExpiration", this.tokenExpiration)

		this.scheduleRefresh()
	}

	scheduleRefresh() {
		this.refreshTimeout = window.setTimeout(this.refreshSession, this.tokenExpiration.getTime() - 1000)
	}

	logout() {
		this.username = ""
		this.token = ""
		this.refreshToken = ""
		sessionStorage.setItem("token", null)
		sessionStorage.setItem("username", null)
		sessionStorage.setItem("tokenExpiration", null)

		window.clearTimeout(this.refreshTimeout)
	}

	computeExpiration() {
		const b64 = this.token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')
		const claims = JSON.parse(decodeURIComponent(atob(b64)))
		return new Date(claims.exp * 1000)
	}
}
