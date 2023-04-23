import { writable } from "svelte/store";

export const view = writable()

export const user_id = writable();
export const session_object = writable();

// export const current_album = writable(0);
// export const current_radio = writable("");
export const album_list = writable({
    albums : [],
    playlists: [],
    radios: []
});

export const active_view = writable({
    name: "",
    id: "",
})

export const track_metadata_view = writable(null);