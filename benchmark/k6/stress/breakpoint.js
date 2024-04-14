import { okLogin } from '../common.js'

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
		tryToBreak: {
			executor: 'ramping-arrival-rate',
			preAllocatedVUs: 500,
			timeUnit: '1s',
			startRate: 0,
			stages: [
				{ target: 10000, duration: '2h' },
			],
		}
	}
};

export default okLogin
