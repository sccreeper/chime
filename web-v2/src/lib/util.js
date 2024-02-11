/**
 * Convert into  hh:mm
 * @param {number} time 
 * @returns {string}
 */
export function convertDurationLong(time) {
    
    let hours = Math.floor(time / 3600)
    let minutes = Math.floor((time % 3600) / 60)

    if (hours == 0) {
        return `${minutes.toString()} min`
    } else {
        return `${hours.toString()} hr ${minutes.toString()} min`
    }

}

// Convert into mm:ss
/**
 * 
 * @param {number} time 
 * @returns {string}
 */
export function convertDuration(time) {
    
    let minutes = Math.floor(time / 60)
    let seconds = Math.floor(time - (minutes * 60)).toString()

    return `${minutes.toString().padStart(2, "0")}:${seconds.padStart(2, "0")}`

}