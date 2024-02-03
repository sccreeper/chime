import { sessionExists } from '$lib/api/auth'
import { redirect } from '@sveltejs/kit'

/** @type {import('./$types').LayoutServerLoad} */
export async function load({cookies}) {
    
    const sessionId = cookies.get("session_id")?.toString()

    if (sessionId === undefined ) {
        throw redirect(303, "/app/login")
    } else if (!await sessionExists(sessionId)) {
        throw redirect(303, "/app/login")
    }

}