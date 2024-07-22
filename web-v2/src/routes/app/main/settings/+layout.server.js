import { getUser } from "$lib/api/auth";
import { redirect } from "@sveltejs/kit";

/** @type {import("./$types").LayoutServerLoad} */
export async function load(cookies) {
    
    const sessionId = cookies.cookies.get("session_id")

    if (sessionId == undefined) {
        throw redirect(303, "/app/login")
    }

    const user = await getUser(sessionId)

    return {
        user: user
    }

}