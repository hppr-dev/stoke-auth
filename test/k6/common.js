import { check, sleep } from 'k6'
import http from 'k6/http'

export const okLogin = function(user, pass, sleepTime) {
  const resp = http.post('http://localhost:8080/api/login',
		JSON.stringify({
			"username" : user,
			"password" : pass,
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

	sleep(sleepTime)

	return resp
}

export const badLogin = function(user, pass, sleepTime) {
  const resp = http.post('http://localhost:8080/api/login',
		JSON.stringify({
			"username" : user,
			"password" : pass,
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

	sleep(sleepTime)

	return resp
}
