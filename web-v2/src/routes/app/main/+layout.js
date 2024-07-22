import { getUserLibrary } from '$lib/api/library'
import { pingServer } from '$lib/api/auth';

/** @type {import('./$types').LayoutLoad} */
export async function load({fetch}) {

    const userLibrary = await getUserLibrary(fetch);

    return {
        lib: userLibrary,
        serverId: await pingServer()
    }

}