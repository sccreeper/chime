import { browser } from "$app/environment"
import { error } from "@sveltejs/kit"
import { createHeaders } from "./api"

/**
 * Get's user's library
 * Server side only
 * @param {import("@sveltejs/kit").Cookies} cookies
 * @returns {Promise<{albums: {id: string, name: string}[], playlists: {id: string, name: string}[], radios: {id: string, name: string}[]}>} userLibrary 
 */
export async function getUserLibrary(cookies) {

    if (cookies.get("user_id") == undefined || cookies.get("session_id") == undefined) {
        throw error(403, "Invalid cookies")
    }

    const req = await fetch(`${browser ? '' : process.env.ORIGIN}/api/get_collections`, {
        // @ts-ignore
        headers: createHeaders(cookies.get("user_id"), cookies.get("session_id")),
    })
    const resp = await req.json()

    return {
        albums: resp.albums,
        playlists: resp.playlists,
        radios: resp.radios,
    }

}