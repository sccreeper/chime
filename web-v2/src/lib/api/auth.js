import { error } from "@sveltejs/kit"
import { browser } from "$app/environment"

/**
 * Checks if a users session is valid and exists
 * Server side only
 * @param {string} sessionId 
 * @returns {Promise<boolean>}
 */
export async function sessionExists(sessionId) {
    
    const req = await fetch(`${browser ? '' : process.env.ORIGIN}/api/auth/session_exists/${sessionId}`)
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
 * @returns {Promise<{session: import("$lib/api/models").Session, successful: boolean}>}
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
        session: {
            sessionId: resp.session.session_id,
            userId: resp.session.user_id,
            isAdmin: resp.user.is_admin,
            username: resp.user.username,
        },
        successful: true,
    }
    
}

// TODO: Merge both of the get user and session exists methods in the future.
/**
 * Gets the user object from a session_exists request. TODO: This will be changed in the future.
 * @param {string} sessionId 
 * @returns {Promise<import("$lib/api/models").Session>}
 */
export async function getUser(sessionId) {
    
    const req = await fetch(`${browser ? '' : process.env.ORIGIN}/api/auth/session_exists/${sessionId}`)
    const resp = await req.json()

    return {
        sessionId: sessionId,
        userId: resp.user.user_id,
        isAdmin: resp.user.is_admin,
        username: resp.user.username,
    }

}

/**
 * Pings server
 * Server responds with "pong" and server id
 * @returns {Promise<string>}
 */
export async function pingServer() {
    
    const req = await fetch(`${browser ? '' : process.env.ORIGIN}/api/ping`)
    const resp = await req.json()

    return resp.server_id

}