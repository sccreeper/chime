// Utility functions

export function convertDuration(time) {
    
    let minutes = Math.floor(time / 60)
    let seconds = Math.floor(time - (minutes * 60)).toString()

    return `${minutes.toString().padStart(2, "0")}:${seconds.padStart(2, "0")}`

}

export function get_url_extension( url ) {
    return url.split(/[#?]/)[0].split('.').pop().trim();
}