import { getStorage } from "$lib/api/admin";

/** 
 * @type {import("./$types").PageLoad} 
 * @returns {Promise<{storageData: import("$lib/api/models").StorageInformation}>}
 */
export async function load({fetch}) {

    return {
        storageData: await getStorage(fetch)
    }

}