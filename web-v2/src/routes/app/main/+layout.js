import { getUserLibrary } from '$lib/api/library'

/** @type {import('./$types').LayoutLoad} */
export async function load({fetch}) {

    let userLibrary = await getUserLibrary(fetch);

    return {
        lib: userLibrary,
    }

}