import { okLogin } from '../../common.js'

export const options = {
	thresholds: {
		http_req_failed: [ // should error out less than 1%
			{
				threshold: 'rate<0.01',
				abortOnFail: true,
				delayAbortEval: '10s'
			}
		],
		http_req_duration: [ // 99% of requestes should be below 0.3s
			{
				threshold: 'p(99)<300',
				abortOnFail: true,
				delayAbortEval: '10s',
			}
		],
	},
	scenarios: {
		nominal_stress: {
			executor: 'ramping-arrival-rate',
			preAllocatedVUs: 50,
			timeUnit: '1s',
			startRate: 0,
			stages: [
				{ target: 10, duration: "15s"},
				{ target: 15, duration: "15s"},
				{ target: 15, duration: '2h' },
				{ target: 10, duration: "30s"},
			],
		}
	}
};

export default okLogin
