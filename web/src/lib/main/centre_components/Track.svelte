<script>
    import { get } from "svelte/store";
    import { active_view, collection_tracks, track_metadata_view } from "../../stores";
    import { convertDuration } from "../../util";
    import FavouriteButton from "./FavouriteButton.svelte";
    import TrackPlay from "./TrackPlay.svelte";

    export let id = "";
    export let title = "";
    export let artist = "";
    export let duration = 0;
    export let album_name = "";
    export let index = 0;
    export let favourited = false;

    let real_index = 0

    $: real_index = index - 1

    function changeView() {
        track_metadata_view.set(id)
    }

    function allowDrop(e) {
        e.preventDefault()
    }


    function drag(e) {
        e.dataTransfer.setData("text", real_index.toString(10))
    }

    function drop(e) {
        e.preventDefault()
        var data = parseInt(e.dataTransfer.getData("text"), 10);

        let track = get(collection_tracks)[data]

        let tracks = get(collection_tracks)
        tracks.splice(data, 1) //Delete moved track from collection
        tracks.splice(real_index, 0, track) //Insert track into collection

        // Send to server

        let track_id_list = []

        tracks.forEach((t) => {
            track_id_list.push(t.id)
        })

        fetch("/api/edit/reorder_collection", {
            method: "POST",
            body: JSON.stringify({
                collection_id: get(active_view).id,
                tracks: track_id_list,
            })
        }).then(resp => {
            if (resp.ok) {
                collection_tracks.set(tracks)      
            } else {
                // TODO: Handle error with toast or something similar.
            }
        })

    }


</script>

<tr on:click={changeView} draggable="true" on:drop={drop} on:dragover={allowDrop} on:dragstart={drag}>
    <td><TrackPlay index={index} track_id={id}/></td>
    <td class="font-semibold">{title}</td>
    <td class="text-xs">{artist}</td>
    <td class="text-xs">{album_name}</td>
    <td class="text-xs">{convertDuration(duration)}</td>
    <td><FavouriteButton id={id} favourited={favourited}/></td>
</tr>

<style>
    td {
        @apply cursor-pointer;
        @apply select-none;
        text-overflow: ellipsis;
        white-space: nowrap;
    }

    tr {
        @apply transition-all;
        @apply text-gray-500;
    }

    tr:hover {
        @apply text-yellow-600;
    }
</style>
