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

var audio = new Audio();
audio.pause();

// Duration change

audio.ontimeupdate = () => {
    position.set(audio.currentTime)
}

position.subscribe(onDurationChange)
function onDurationChange() {
    audio.currentTime = get(position);
}

