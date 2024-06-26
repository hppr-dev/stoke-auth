import { badLogin, okLogin } from './common.js'

export const options = {
	scenarios: {
		issued: {
			executor: "constant-vus",
			exec: "issued",
			vus: 5,
			duration: '10s',
		},
		rejected: {
			executor: "constant-vus",
			exec: "rejected",
			vus: 5,
			duration: '10s',
		},
	}
};

export const issued = () => okLogin("tester", "tester", 0.01);

export const rejected = () => badLogin("tester", "badpass", 0.01);
