import { get, writable } from "svelte/store"
import { duration, playing, position, volume } from "./player"

export const cast_devices = writable([])
export const cast_uuids = writable([])
export const current_cast_device = writable({uuid: "", name: ""})
export const using_cast = writable(false)

 export const CastStates = {
    Stop : "stop",
    Pause : "pause",
    Play : "play"
}

let cast_states = ["stop", "pause", "play"]

// Methods for interacting with the Cast proxy API.
export const Cast = {

    // Get all of the cast devices
    discover: () => {
        fetch("/api/cast/get_devices").then(resp => resp.json()).then(
            data => {
                cast_devices.set(data)

                let uuids = []

                data.forEach(element => {
                    
                    uuids.push(element.uuid)

                });

                cast_uuids.set(uuids)
            }
        )   
    },

    // Set the UUID passed as the current cast device that is being used.
    connect: (uuid) => {

        if (get(cast_uuids).includes(uuid)) {

            using_cast.set(true)

            for (let i = 0; i < get(cast_devices).length; i++) {
                
                if (get(cast_devices)[i].uuid == uuid) {
                    current_cast_device.set({uuid: uuid, name: get(cast_devices)[i].name})
                    break
                }

            }
        } else {
            console.error("Invalid UUID")
        }
        
    },

    // Stop casting
    stop: (uuid) => {

        if (get(cast_uuids).includes(uuid)) {
            using_cast.set(false)
            current_cast_device.set({name: "", uuid: ""})
            
            fetch("/api/cast/control", {
                method: "POST",
                body: JSON.stringify({
                    uuid: uuid,
                    state: "stop"
                })
            })

            using_cast.set(false)
        } else {
            console.error("Invalid UUID")
        }

    },

    // Play media
    play_media: async (uuid, url, mimetype) => {

        if (get(cast_uuids).includes(uuid)) {
            
            await fetch("/api/cast/play_media",  {
                method: "POST",
                body: JSON.stringify({
                    uuid: uuid,
                    url: url,
                    mimetype: mimetype
                })
            })

            playing.set(true)
            setTimeout(cast_update, 1000)

        } else { 
            console.error("Invalid UUID")
        }

    },

    // Get's the current cast status and returns it.
    get_status: async (uuid) => {
        if (get(using_cast) && get(playing) && get(cast_uuids).includes(uuid)) {
        
            let status_req = await fetch(`/api/cast/get_status/${uuid}`, {
                method: "GET"
            })
            let status_json = await status_req.json()
            return status_json

        }

    },

    set_volume: (uuid, volume) => {

        if (get(cast_devices).includes(uuid)) {
            
            fetch("/api/cast/set_volume", {
                method: "POST",
                body: JSON.stringify({
                    uuid: uuid,
                    volume: volume
                })
            })   
        } else {
            console.error("Invalid UUID")
        }

    },

    control: (uuid, state) => {
        
        if (cast_states.includes(state) && get(cast_uuids).includes(uuid)) {
            
            fetch("/api/cast/control", {
                method: "POST",
                body: JSON.stringify({
                    uuid: uuid,
                    state: state
                })
            })

        } else {
            console.error("Invalid UUID")
        }

    }

}

// Async method to fetch cast status and then update the UI.
const cast_update = async () => {

    console.log(get(using_cast))

    if (get(using_cast)) {

        let status = await Cast.get_status(get(current_cast_device).uuid)

        position.set(status.current_time)
        volume.set(status.volume)
        duration.set(status.duration)

        setTimeout(cast_update, 1000)
    }

}