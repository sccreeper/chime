// Utility functions

export function convertDuration(time) {
    
    let minutes = Math.floor(time / 60).toString()
    let seconds = Math.floor((time % 60) * 60).toString()

    return `${minutes.padStart(2, "0")}:${seconds.padStart(2, "0")}`

}