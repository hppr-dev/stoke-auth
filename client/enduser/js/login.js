const api_url = "http://localhost:8080"
var token = "",
	refreshToken = "",
	username = "",
	tokenExpiration = Date.now()

async function login(user, password, required_claims, callback) {
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

	callback()

	username = user
	token = result.token
	refreshToken = result.refresh
	tokenExpiration = computeExpiration(token)

	sessionStorage.setItem("token", result.token)
	sessionStorage.setItem("refresh", result.refresh)
	sessionStorage.setItem("username", username)

	scheduleRefresh()
}

async function refreshSession() {
  const response = await fetch(`${api_url}/api/refresh`, {
    method: "POST",
    headers: {
      "Content-Type" : "application/json",
      "Authorization" : `Bearer ${token}`,
    },
    body: JSON.stringify({
      refresh: refreshToken,
    }),
  })

  if ( !response.ok ){
    throw new Error(response.statusText)
  }

  const result = await response.json();
  if ( result.message ) {
    throw new Error(result.message);
  }

  token = result.token
  refreshToken = result.refresh
  sessionStorage.setItem("token", result.token)
  sessionStorage.setItem("refresh", result.refresh)

  scheduleRefresh()
}

function scheduleRefresh() {
	refreshTimeout = window.setTimeout(refreshSession, tokenExpiration.getTime() - Date.now() - 10000)
}

function logout() {
	username = ""
	token = ""
	refreshToken = ""
	sessionStorage.setItem("token", "")
	sessionStorage.setItem("refresh", "")
	sessionStorage.setItem("username", "")
	window.clearTimeout(refreshTimeout)
}

function computeExpiration(token) {
  const b64 = token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')
  const claims = JSON.parse(decodeURIComponent(atob(b64)))
  return new Date(claims.exp * 1000)
}

export logout, login
