import { badLogin, okLogin } from './common.js'

export const options = {
	scenarios: {
		issued: {
			executor: "constant-arrival-rate",
			exec: "issued",
			duration: '30m',
			rate: 50,
			preallocatedVUs: 25,
		},
		rejected: {
			executor: "constant-arrival-rate",
			exec: "rejected",
			duration: '30m',
			rate: 50,
			preallocatedVUs: 25,
		},
	}
};

export const issued = () => okLogin("tester", "tester", 0.01);

export const rejected = () => badLogin("tester", "badpass", 0.01);
