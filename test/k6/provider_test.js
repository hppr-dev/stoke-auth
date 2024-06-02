import http from 'k6/http'
import { badLogin, okLogin } from './common.js'

export const options = {
	scenarios: {
		providerIssuedUser: {
			executor: "constant-vus",
			exec: "providerIssuedUser",
			vus: 2,
			duration: '10s',
		},
		localIssued: {
			executor: "constant-vus",
			exec: "localIssued",
			vus: 2,
			duration: '10s',
		},
		providerIssuedAdmin: {
			executor: "constant-vus",
			exec: "providerIssuedAdmin",
			vus: 2,
			duration: '10s',
		},
		localRejected: {
			executor: "constant-vus",
			exec: "localRejected",
			vus: 2,
			duration: '10s',
		},
		providerRejected: {
			executor: "constant-vus",
			exec: "providerRejected",
			vus: 2,
			duration: '10s',
		},
	}
};

export function setup() {
	const requests = [
		JSON.stringify({ "username" : "leela", "password": "leela" }),
		JSON.stringify({ "username" : "fry", "password": "fry" }),
		JSON.stringify({ "username" : "hermes", "password": "hermes" }),
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

export const providerIssuedAdmin = () => okLogin("fry", "fry", 0.01);

export const providerIssuedUser = () => okLogin("hermes", "hermes", 0.01);

export const providerRejected = () => badLogin("leela", "badpass", 0.01);

export const localIssued = () => okLogin("tester", "tester", 0.01);

export const localRejected = () => badLogin("tester", "badpass", 0.01);
