import { browser } from "$app/environment"

/**
 * Get's user's library
 * Server side only
 * @param {fetch} fetch
 * @returns {Promise<{albums: {id: string, name: string}[], playlists: {id: string, name: string}[], radios: {id: string, name: string}[]}>} userLibrary 
 */
export async function getUserLibrary(fetch) {

    const req = await fetch(`${browser ? '' : process.env.ORIGIN}/api/get_collections`)
    const resp = await req.json()

    return {
        albums: resp.albums,
        playlists: resp.playlists,
        radios: resp.radios,
    }

}

/**
 * 
 * @param {fetch} fetch 
 * @param {string} id 
 * 
 * @returns {Promise<import("./api").Collection>}
 */
export async function getCollection(fetch, id) {
    
    const req = await fetch(`${browser ? '' : process.env.ORIGIN}/api/get_collection/${id}`)
    const collection = await req.json()

    return collection

}