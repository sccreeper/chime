import { get } from "svelte/store";
import { audio_source, player_audio, playing, repeat, shuffle, volume } from "./player";
import { clamp } from "./util";

export function handle_keydown(e) {
    switch (e.code) {
        // Change volume
        case "ArrowDown":
            volume.set(clamp(get(volume)-0.05, 0, 1))
            break;
        case "ArrowUp":
            volume.set(clamp(get(volume)+0.05, 0, 1))
            break;
        
        // Skip forwards/backwareds
        case "ArrowRight":
            if (get(audio_source).type == "track" && get(playing)) {
                player_audio.currentTime+=10;
            }
            break;
        case "ArrowLeft":
            if (get(audio_source).type == "track" && get(playing)) {
                player_audio.currentTime-=10;
            }
            break;

        // Toggle shuffle and repeat
        case "KeyR":
            console.log("repeat")
            repeat.set(!get(repeat))
            break;
        case "KeyS":
            shuffle.set(!get(shuffle))
            break;
        
        // Play/pause
        case "Space":
            playing.set(!get(playing))
            break;
        default:
            break;
    }
}