import { browser } from "$app/environment"

/**
 * @typedef {Object} LibraryItem
 * @property {string} id
 * @property {string} name
 * @property {string} cover_id
 */

/**
 * Get's user's library
 * Server side only
 * @param {fetch} fetch
 * @returns {Promise<{albums: LibraryItem[], playlists: LibraryItem[], radios: LibraryItem[]}>} userLibrary 
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
 * @returns {Promise<import('./api').Collection>}
 */
export async function getCollection(fetch, id) {
    
    const req = await fetch(`${browser ? '' : process.env.ORIGIN}/api/get_collection/${id}`)
    const collection = await req.json()

    return collection

}

/**
 * 
 * @param {string} id 
 * @returns {Promise<import('./api').TrackMetadata>}
 */
export async function getTrackMetadata(id) {
    
    const req = await fetch(`/api/get_track_metadata/${id}`)
    const metadata = await req.json()

    return metadata

}