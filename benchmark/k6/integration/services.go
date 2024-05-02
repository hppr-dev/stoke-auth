import http from 'k6/http'
import { check } from 'k6'

export const options = {
  vus: 5,
  duration: '1m',
};

// Setup:
// * stoke server on 8080 with ldap integration
// * client/examples/go/engine on 5000
export default function() {
	const services = [ 
	"http://localhost:5000/speed", // requires role:eng
	// TODO add more 
	]
  const stokeResp = http.post('http://localhost:8080/api/login',
		JSON.stringify({
			"username" : "leela",
			"password" : "leela",
		}),
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
	})

	return stokeResp
}
