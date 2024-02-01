import { loginUser } from '$lib/api/auth';
import { error, redirect } from '@sveltejs/kit';

const MAX_COOKIE_AGE = 7 * 24 * 3600

/** @type {import('./$types').Actions} */
export const actions = {
    default: async ({cookies, request}) => {
        
        const form = await request.formData();
        const username = form.get("username")?.toString()
        const password = form.get("password")?.toString()

        if (username === undefined || password === undefined) {
            throw error(400, "Bad request")
        }
        const auth = await loginUser(username, password)

        if (auth.successful) {
            cookies.set("session_id", auth.sessionId, {path: "/", sameSite: false, httpOnly: false, secure: false, maxAge: MAX_COOKIE_AGE})
            cookies.set("user_id", auth.userId, {path: "/", sameSite: false, httpOnly: false, secure: false, maxAge: MAX_COOKIE_AGE})
            cookies.set("is_admin", auth.isAdmin, {path: "/", sameSite: false, httpOnly: false, secure: false, maxAge: MAX_COOKIE_AGE})
    
            throw redirect(303, "/app/main")   
        } else {
            throw redirect(303, "/app/login")   
        }

    }
}