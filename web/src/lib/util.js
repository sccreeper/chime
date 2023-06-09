// Utility functions & constants

export const allowed_username_chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890_-"

export function convertDuration(time) {
    
    let minutes = Math.floor(time / 60)
    let seconds = Math.floor(time - (minutes * 60)).toString()

    return `${minutes.toString().padStart(2, "0")}:${seconds.padStart(2, "0")}`

}

export function getUrlExtension( url ) {
    return url.split(/[#?]/)[0].split('.').pop().trim();
}

export function verifyString(s, check) {
    
    for(let char of s) {

        if (!check.includes(char)) {
            return false
        }

    }

    return true

}

export function clamp(val, min, max) {
    return Math.max(min, Math.min(val, max))
}