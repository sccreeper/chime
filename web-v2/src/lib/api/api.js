export const serverPort = "8042"

/**
 * 
 * @param {string} user_id 
 * @param {string} session_id 
 * @returns {Headers}
 */
export function createHeaders(user_id, session_id) {
    
    let h = new Headers();
    h.set("Cookie", `user_id=${user_id}; session_id=${session_id}`)
    return h

}