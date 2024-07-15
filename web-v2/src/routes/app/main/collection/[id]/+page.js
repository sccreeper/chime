import { getCollection } from '$lib/api/library';
import { error } from '@sveltejs/kit';

/**  
 * @type {import('./$types').PageLoad}
 * @returns {Promise<{collection: import('$lib/api/models').Collection, collection_id: string}>}
 */
export async function load({params, fetch}) {
    
    if (params.id.length == 0) {
        throw error(400, "Invalid ID")
    }

    const collection = await getCollection(fetch, params.id)

    return {
        collection: collection,
        collection_id: params.id
    }

}