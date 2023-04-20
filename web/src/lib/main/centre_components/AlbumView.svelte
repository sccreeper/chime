<script>
    import { current_album, session_object } from "../../stores";
    import { get } from "svelte/store";
    import Track from "./Track.svelte";


let album_title = ""
let album_cover_src = ""
let tracks = []

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

<div class="inline">
    <img src="" width="256" height="256"/>
    <h1>{album_title}</h1>
</div>

<div>

{#if tracks.length == 0}

<p>No tracks in album.</p>

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

