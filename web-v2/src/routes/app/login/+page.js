import { pingServer } from '$lib/api/auth';

/** @type {import('./$types').PageLoad} */
export async function load() {
	return {
        server_id: await pingServer(),
	};
}