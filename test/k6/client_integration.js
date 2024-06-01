import http from 'k6/http'
import { check, sleep } from 'k6'

export const options = {
  vus: 10,
  duration: '1m',
};

export function setup() {
	const requests = [
		JSON.stringify({ "username" : "leela", "password": "leela" }),
		JSON.stringify({ "username" : "fry", "password": "fry" }),
	]

	// Make request for users to create them in the database
	requests.forEach((rq) => http.post('http://localhost:8080/api/login',
			rq,
			{
				headers: {
					"Content-Type" : "application/json"
				}
			}
		)
	)
}

// Setup:
// * stoke server on 8080 with ldap integration
// * docker compose file client/client-test-compose.yaml running
export default function() {
	const services = [ 
	"http://localhost:8888/control/location",  // requires ctl:nav
	"http://localhost:8888/request/shipment",  // requires req:acc
	// TODO add more 
	]
	const requests = [
		JSON.stringify({ "username" : "leela", "password": "leela" }),
		JSON.stringify({ "username" : "fry", "password": "fry" }),
	]
  const stokeResp = http.post('http://localhost:8080/api/login',
		requests[Math.floor(Math.random() * 2)],
		{
			headers: {
				"Content-Type" : "application/json"
			}
		}
	)

	check(stokeResp, {
		"token response code was 200": (resp) => resp.status == 200,
	})

	const stokeBody = JSON.parse(stokeResp.body)
	check(stokeBody, {
		"response contained token": (body) => body.token && body.refresh,
	})

	services.forEach((service) => {
		let checks = {}
		checks[`${service} response 200`] = (resp) => resp.status == 200
		const resp = http.get(service, {
			headers: {
				"Authorization": "Bearer " + stokeBody.token
			}
		})
		check(resp, checks)
		sleep(Math.random() * 4)
	})

	return stokeResp
}
