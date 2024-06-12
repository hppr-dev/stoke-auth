import http from 'k6/http'
import ws from 'k6/ws'
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
	// http services
	const services = [ 
		"http://localhost:8888/control/location",            // requires ctl:nav -- go rest
		"http://localhost:8888/control/speed",               // requires ctl:sp  -- go rest/unary grpc
		"http://localhost:8888/request/shipment",            // requires req:acc -- python rest flask
		"http://localhost:8888/inventory/test/",             // requires inv:acc -- python rest django
		"http://localhost:8888/inventory/cargo_contents/",   // requires car:acc -- python rest django/unary grpc
	]
	//ws services. Tokens are sent as url parameters
	const ws_services = [
		{ url: "ws://localhost:8888/control/foobar", request: "foo", response: "bar", times: 3 },    // requires ctl:acc  -- go rest/stream grpc
	]
	const user_logins = [
		JSON.stringify({ "username" : "leela", "password": "leela" }),
		JSON.stringify({ "username" : "fry", "password": "fry" }),
	]
  const stokeResp = http.post('http://localhost:8080/api/login',
		user_logins[Math.floor(Math.random() * 2)],
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
		if (resp.status != 200) {
			console.log("Request failed!",service, resp)
		}
		sleep(Math.random() * 4)
	})

	ws_services.forEach((service) => {
		let checks = {}
		checks[`${service.url} send ${service.request} -> recv ${service.response}`] = (data) => data == service.response
		const res = ws.connect(service.url + "?token=" + stokeBody.token, {}, function (socket) {
			let times = 0
			socket.on('open', () => {
				socket.send(service.request)
				times += 1
			});
			socket.on('message', (data) => {
				check(data, checks)
				socket.send(service.request)
				times += 1
				if (service.times == times) {
					socket.close()
				}
			});
		});

		check(res, { 'ws response is 101': (r) => r && r.status === 101 });
	})

	return stokeResp
}
