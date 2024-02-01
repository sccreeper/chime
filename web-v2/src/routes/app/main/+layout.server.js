import { sessionExists } from '$lib/api/auth'
import { getUserLibrary } from '$lib/api/library'
import { redirect } from '@sveltejs/kit'

/** @type {import('./$types').LayoutServerLoad} */
export async function load({ cookies, request }) {
    
    const sessionId = cookies.get("session_id")?.toString()

    if (sessionId === undefined ) {
        throw redirect(303, "/app/login")
    } else if (!await sessionExists(sessionId)) {
        throw redirect(303, "/app/login")
    }

    let userLibrary = await getUserLibrary(cookies);

    console.log(userLibrary)

    return {
        lib: userLibrary,
    }

}