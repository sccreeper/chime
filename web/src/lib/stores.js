// General stores used across application.
import { writable } from "svelte/store";

export const view = writable()

export const session_object = writable();
export const user_object = writable({username: "", is_admin: false, user_id: ""});

export const album_list = writable({
    albums : [],
    playlists: [],
    radios: []
});
export const collection_tracks = writable([]);

export const active_view = writable({
    name: "",
    id: "",
})

export const track_metadata_view = writable(null);

export const favourites_list = writable([])

export const search_results = writable({
    tracks: [],
    collections: [],
    radios: [],
});

export const edit_user_error = writable("");