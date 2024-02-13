import { browser } from "$app/environment";
import { get, writable } from "svelte/store";
import { convertDuration } from "./util";

export const PLAYER_CONTEXT_KEY = "chime_player"
export const CURRENT_COLLECTION_CONTEXT_KEY = "current_collection"

export class ChimePlayer {

    // Player state

    /** @type {import('svelte/store').Writable<number>} */
    duration = writable(0);
    currentTime = writable(0);

    /** @type {import('svelte/store').Writable<string>} */
    durationString = writable("");
    currentTimeString = writable("");

    /** @type {import('svelte/store').Writable<boolean>} */
    playing = writable(false);

    // Track information

    /** @type {import('svelte/store').Writable<string>} */
    collectionId = writable("");
    /** @type {import('svelte/store').Writable<import('$lib/api/api').Collection | undefined>} */
    currentCollection = writable();
    /** @type {import('svelte/store').Writable<import('$lib/api/api').CollectionTrack | undefined>} */
    currentTrack = writable();

    /** @type {number} */
    collectionIndex = 0;

    /** @type {HTMLAudioElement | null} */
    audioElement = null;

    constructor () {
        
        // If in browser (client) environment, init audio player.

        if (browser) {
            
            this.audioElement = new Audio();
            this.playing.set(false);

            this.audioElement.addEventListener("timeupdate", () => {
                if (this.audioElement?.currentTime == undefined) {
                    this.currentTime.set(0)   
                } else {
                    this.currentTime.set(this.audioElement.currentTime)
                }
            })

            this.audioElement.addEventListener("canplaythrough", () => {
                if (this.audioElement?.duration == undefined) {
                    this.duration.set(0)
                } else {
                    this.duration.set(this.audioElement.duration)
                }
            })

            // Store subscriptions

            this.duration.subscribe((val) => {
                this.durationString.set(convertDuration(val))
            })

            this.currentTime.subscribe((val) => {
                this.currentTimeString.set(convertDuration(val))
            })

            this.playing.subscribe((state) => {
                if (this.audioElement == undefined) {
                    return
                } else {

                    if (state) {
                        this.audioElement.play()   
                    } else {
                        this.audioElement.pause()
                    }
                }
            })

        }

    }

    /**
     * Plays a track
     * @param {import('$lib/api/api').CollectionTrack} track 
     */
    #playTrack(track) {

        if (this.audioElement == null) {
            return
        }

        this.audioElement.src = `/api/stream/${track.id}`
        this.currentTrack.set(track)
        this.playing.set(true)

    }

    /**
     * Plays a collection, with given starting index.
     * @param {import('$lib/api/api').Collection} collection 
     * @param {number} index 
     * @param {string} id
     */
    playCollection(collection, index, id) {

        if (this.audioElement == null) {
            return
        }

        this.audioElement.pause()
        this.playing.set(false)

        this.currentCollection.set(collection)
        this.collectionId.set(id)
        this.collectionIndex = index

        this.#playTrack(collection.tracks[index])

    }

    /**
     * Seek to position on track
     * @param {number} time 
     */
    seek(time) {

        if (this.audioElement == undefined) {
            return
        } else {
            this.audioElement.currentTime = time
        }
    }


}