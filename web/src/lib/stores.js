import { SvelteComponent } from "svelte/internal";
import { writable } from "svelte/store";

export const view = writable()

export const user_id = writable();
export const session_object = writable();

export const current_album = writable(0);
export const album_list = writable({
    albums : [],
    playlists: [],
});