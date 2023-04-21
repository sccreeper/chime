<script>
    import { current_album, session_object } from "../../stores";
    import { get } from "svelte/store";
    import Track from "./Track.svelte";
    import no_cover_image from "../../../assets/no_cover.png"


let album_title = ""
let album_cover_src = ""
let tracks = []
let is_album = false

function updateView(album_id) {
    
    fetch("/api/get_collection", {
        method: "POST",
        body: JSON.stringify({
            session : get(session_object),
            album_id: album_id,
        }),
    }).then(response => response.json()).then(data => {

        album_title = data.title
        tracks = data.tracks


    })

}

current_album.subscribe(value => updateView(value));

</script>

<div class="m-2">

<div class="flex flex-row items-center gap-4">
    <img src={album_cover_src == "" ? no_cover_image : album_cover_src} class="album-cover"/>
    <h1 class="album-title float-right">{album_title}</h1>
</div>

<hr class="m-3 h-px border-none bg-gray-600">

<div>

{#if tracks.length == 0}

<p>No tracks in {is_album ? "album" : "playlist"}. Add from the lefthand sidebar or upload files.</p>

{:else}

<table>

    <tr>
        <th>Title</th>
        <th>Artist</th>
        <th>Album</th>
        <th>Duration</th>
    </tr>

</table> 

{#each tracks as track}
    
    <Track id={track.id} 
    title={track.name} 
    artist={track.artist} 
    duration={track.duration} 
    album_name={track.album_name}
    />

{/each}

{/if}


</div>

</div>

<style>

    .album-title { 

        @apply text-6xl;
        @apply text-gray-300;
        @apply font-bold;

    }

    .album-cover {

        width: 256px;
        height: 256px;

        min-height: 128px;
        min-width: 128px;

    }

</style>
