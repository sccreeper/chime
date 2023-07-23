import { get, writable } from "svelte/store"

export const cast_devices = writable([])
export const cast_uuids = writable([])
export const current_cast_device = writable({uuid: "", name: ""})
export const using_cast = writable(false)

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

        using_cast.set(true)

        if (get(cast_uuids).includes(uuid)) {

            for (let i = 0; i < get(cast_devices).length; i++) {
                
                if (get(cast_devices)[i].uuid == uuid) {
                    current_cast_device.set({uuid: uuid, name: get(cast_devices)[i].name})
                    break
                }

            }
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
        }

    },

    // Play media
    play_media: (uuid, url, mimetype) => {

        if (get(cast_uuids).includes(uuid)) {
            
            fetch("/api/cast/play_media",  {
                method: "POST",
                body: JSON.stringify({
                    uuid: uuid,
                    url: url,
                    mimetype: mimetype
                })
            })

        }

    }

}