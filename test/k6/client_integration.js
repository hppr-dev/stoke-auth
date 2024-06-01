import http from 'k6/http'
import { check, sleep } from 'k6'

export const options = {
  vus: 20,
  duration: '1m',
};

// Setup:
// * stoke server on 8080 with ldap integration
// * docker compose file client/client-test-compose.yaml running
export default function() {
	const services = [ 
	"http://localhost:5000/location",          // requires ctl:nav
	"http://localhost:5001/shipment-request",  // requires req:acc
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
		const resp = http.get(service, {
			headers: {
				"Authorization": "Bearer " + stokeBody.token
			}
		})
		check(resp, {
			"service response code was 200": (resp) => resp.status == 200,
		})
		sleep(Math.random() * 4)
	})

	return stokeResp
}
