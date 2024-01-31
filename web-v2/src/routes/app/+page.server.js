import { sessionExists } from '$lib/api/auth'
import { redirect } from '@sveltejs/kit'

/** @type {import('./$types').PageServerLoad} */
export async function load({ cookies }) {
 
    const sessionId = cookies.get("session_id")

    if (sessionId === undefined ) {
        throw redirect(307, "/app/login")
    } else if (await sessionExists(sessionId)) {
        throw redirect(307, "/app/main")
    }

}