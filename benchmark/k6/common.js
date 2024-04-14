import http from 'k6/http'
import { check } from 'k6'

export const okLogin = function() {
  const resp = http.post('https://localhost:8080/api/login',
		JSON.stringify({
			"username" : "tester",
			"password" : "tester",
		}),
		{
			headers: {
				"Content-Type" : "application/json"
			}
		}
	)

	check(resp, {
		"response code was 200": (resp) => resp.status == 200,
	})

	const token = JSON.parse(resp.body)
	check(token, {
		"response contained token": (token) => token.token && token.refresh,
	})

	return resp
}

export const badLogin = function() {
  const resp = http.post('https://localhost:8080/api/login',
		JSON.stringify({
			"username" : "tester",
			"password" : "badpass",
		}),
		{
			headers: {
				"Content-Type" : "application/json"
			}
		}
	)

	check(resp, {
		"response code was 401": (resp) => resp.status == 401,
	})

	const token = JSON.parse(resp.body)
	check(token, {
		"response did not contain a token": (token) => token.token === undefined && token.refresh === undefined,
	})

	return resp
}
