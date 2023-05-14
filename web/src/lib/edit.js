// Contains API calls for editing track details

import { favourites_list, session_object } from "./stores"

export function api_favourite_track(track_id) {
    
    fetch("/api/edit/favourite", {
        method: "POST",
        body: JSON.stringify({track_id: track_id})
    }).then(resp => (resp.json()).then(
        data => {
            favourites_list.set(data.favourites)
        }
    ))

}

export function api_get_favorites() {
    fetch("/api/user/get_favourites", {
        method: "GET",
    }).then(resp => (resp.json()).then(
        data => {
            favourites_list.set(data.favourites)
        }
    ))
}

session_object.subscribe(() => {
    api_get_favorites()
})