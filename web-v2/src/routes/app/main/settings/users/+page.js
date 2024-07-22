import { getUsers } from "$lib/api/admin";

/** @type {import("./$types").PageLoad} */
export async function load({fetch}) {
    
    return {
        users: await getUsers(fetch)
    }

}