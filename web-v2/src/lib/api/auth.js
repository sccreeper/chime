import { error } from "@sveltejs/kit"
import { browser } from "$app/environment"

/**
 * Checks if a users session is valid and exists
 * Server side only
 * @param {string} id 
 * @returns {Promise<boolean>}
 */
export async function sessionExists(id) {
    
    const req = await fetch(`${browser ? '' : process.env.ORIGIN}/api/auth/session_exists/${id}`)
    if (req.status != 200) {
        return false
    }

    let body = await req.json()

    if (body.status == "exists") {
        return true
    } else {
        return false
    }

}

/**
 * Authenticates a user
 * Server side only
 * @param {string} username 
 * @param {string} password 
 * @returns 
 */
export async function loginUser(username, password) {

    let f = new FormData()
    f.set("u", username)
    f.set("p", password)

    const req = await fetch(`${browser ? '' : process.env.ORIGIN}/api/auth`, {
        body: f,
        method: "POST"
    })

    if (!req.ok) {
        throw error(500, "Unable to authenticate")
    }

    const resp = await req.json()

    return {
        sessionId: resp.session.session_id,
        userId: resp.session.user_id,
        isAdmin: resp.user.is_admin,
        successful: true,
    }
    
}