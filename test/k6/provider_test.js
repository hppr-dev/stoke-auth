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

export const providerIssuedAdmin = () => okLogin("fry", "fry", 0.01);

export const providerIssuedUser = () => okLogin("hermes", "hermes", 0.01);

export const providerRejected = () => badLogin("leela", "badpass", 0.01);

export const localIssued = () => okLogin("tester", "tester", 0.01);

export const localRejected = () => badLogin("tester", "badpass", 0.01);
