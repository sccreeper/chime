//@ts-nocheck
// Source code for managing the player state.

import { get, writable } from "svelte/store";
import { track_metadata_view } from "./stores";
import { get_url_extension } from "./util";
import "https://cdn.jsdelivr.net/npm/hls.js";

export var playing = writable(false)
export var playing_radio = writable(false)
export var audio_source = writable({type: "", source: ""})
export var position = writable(0)
export var duration = writable(0)
export var track_queue = writable([])

// Playback settings

export var shuffle = writable(false)
export var volume = writable(1.0)

let playing_hls = false;

// HLS player

var hls_player = document.createElement("video")
hls_player.pause()
var hls = new Hls();
hls.attachMedia(hls_player)

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

playing.subscribe((val) => {

    if (val) {

        if (playing_hls) {
            hls_player.play()
        } else {
            player_audio.play()
        }
        
    } else {

        if (playing_hls) {
            hls_player.pause()
        } else {
            player_audio.pause()
        }
    }

})

player_audio.addEventListener("play", () => {
    if (!playing_hls) {
        playing.set(true)   
    }
})

player_audio.addEventListener("pause", () => {

    if (!playing_hls) {
        playing.set(false)   
    }
})

hls_player.addEventListener("play", () => {
    if (playing_hls) {
        playing.set(true)   
    }
})

hls_player.addEventListener("pause", () => {
    if (playing_hls) {
        playing.set(false)   
    }
})

// Volume events

volume.subscribe(() => {
    player_audio.volume = get(volume)
    hls_player.volume = get(volume)
})

// Track change events.

audio_source.subscribe((val) => {

    if (val != null) {

        if (val.type == "track") {

            hls.stopLoad()
            hls_player.pause()

            track_metadata_view.set(val.source)

            player_audio.src = `/api/stream/${val.source}`
            player_audio.play()
            
            playing_radio.set(false)
            playing_hls = false
        } else if (val.type == "radio") {

            if (get_url_extension(val.source) == "m3u8") {

                player_audio.src = ""
                player_audio.pause()

                if (Hls.isSupported()) {
                    console.log("Playing with HLS")

                    hls.loadSource(val.source);
                    hls.startLoad()
                    hls_player.play()
    
                    playing_radio.set(true)
                    playing_hls = true
                }

            } else {

                hls.stopLoad()
                hls_player.pause()

                player_audio.src = val.source
                player_audio.play()

                playing_radio.set(true)
                playing_hls = false
            }
        }
    }

})

