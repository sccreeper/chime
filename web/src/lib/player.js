// Source code for managing the player state.

import { get, writable } from "svelte/store";

export var playing = writable(false)
export var track_id = writable(null)
export var position = writable(0)
export var duration = writable(0)
export var track_queue = writable([])

// Playback settings

export var shuffle = writable(false)
export var volume = writable(1.0)

// AudioPlayer

export var player_audio = new Audio();
player_audio.pause();

// Duration change

player_audio.addEventListener("timeupdate", () => {
    position.set(player_audio.currentTime)
})

player_audio.addEventListener("canplaythrough", () => {
    duration.set(player_audio.duration)
})

// Playing events

playing.subscribe(() => {

    if (get(playing)) {
        player_audio.play()
        playing.set(true)
    } else {
        player_audio.pause()
        playing.set(false)
    }

})

player_audio.addEventListener("play", () => {
    playing.set(true)
})

player_audio.addEventListener("pause", () => {
    playing.set(false)
})

// Volume events

volume.subscribe(() => {
    player_audio.volume = get(volume)
})

// Track change events.

track_id.subscribe(() => {

    if (get(track_id) != null) {
        player_audio.src = `/api/stream/${get(track_id)}`
        player_audio.play()
    }

})

