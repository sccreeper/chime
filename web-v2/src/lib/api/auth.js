/**
 * 
 * @param {string} id 
 * @returns {Promise<boolean>}
 */
export async function sessionExists(id) {
    
    const req = await fetch(`/api/auth/session_exists/${id}`)
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