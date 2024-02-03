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