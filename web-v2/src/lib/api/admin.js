/**
 * 
 * @param {typeof fetch} fetch
 * @returns {Promise<import("$lib/api/models").StorageInformation>}
 */
export async function getStorage(fetch) {

    const storageReq = await fetch("/api/admin/storage")
    const storageData = await storageReq.json()
    
    return {
        totalVolumeSpace: storageData.total_volume_space,
        usedByOthers: storageData.used_by_others,
        usedByChime: storageData.used_by_chime,
        breakdown: {
            backups: storageData.breakdown.backups,
            cache: storageData.breakdown.cache,

            covers: storageData.breakdown.covers,
            tracks: storageData.breakdown.tracks,
        }
    }

}

/**
 * 
 * @param {typeof fetch} fetch 
 * @returns {Promise<import("$lib/api/models").User[]>}
 */
export async function getUsers(fetch) {
    
    const req = await fetch(
        "/api/admin/users",
        {
            method: "GET"
        }
    )

    /** @type {{users: any[]}} */
    const data = await req.json()

    /** @type {import("$lib/api/models").User[]} */
    let userList = []

    data.users.forEach((user) => {
        
        userList.push({
            username: user.username,
            id: user.id,
            isAdmin: user.is_admin,
        })

    });

    return userList

}