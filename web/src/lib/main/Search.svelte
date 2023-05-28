<script>
    import {search_results} from "../stores";
    import CollectionCard from "./search_components/CollectionCard.svelte";
    import TrackResult from "./search_components/TrackResult.svelte";

    var search_playlists = null
    var search_albums = null

    // Subscribe to split collections value when search results are updated.
    search_results.subscribe((val) => {

        search_playlists = null
        search_albums = null

        if (val.collections != null) {
            val.collections.forEach(element => {
                if (element.is_album) {
                    
                    if (search_albums == null) {
                        search_albums = []
                    }

                    search_albums.push(element)

                } else {

                    if (search_playlists == null) {
                        search_playlists = []
                    }

                    search_playlists.push(element)

                }
            });   
        }

    })


</script>

<div class="w-full h-full m-3 overflow-y-scroll">

{#if $search_results.tracks != null}

<h1>Tracks</h1>

<table>

<colgroup>
    <col span="1" style="width: 10%;">
    <col span="1" style="width: 40%;">
    <col span="1" style="width: 35%;">
    <col span="1" style="width: 15%;">
</colgroup>

<tr style="text-align: left;">
    <th><i class="bi bi-image text-gray-500"></i></th>
    <th>Title</th>
    <th>Artist</th>
    <th>Duration</th>
</tr>

{#each $search_results.tracks as track}
    <TrackResult album_id={track.album_id} title={track.title} cover={track.cover} duration={track.duration} artist={track.artist}/>
{/each}

</table>

{/if}

{#if search_albums != null}
<h1>Albums</h1>

<div class="flex flex-row">
{#each search_albums as album}
    <CollectionCard type="Album" view_type="album" id={album.id} name={album.title} cover={album.cover}/>
{/each}
</div>

{/if}

{#if search_playlists != null}
<h1>Playlists</h1>

<div class="flex flex-row">
{#each search_playlists as playlist}
    <CollectionCard type="Playlist" view_type="album" id={playlist.id} name={playlist.title} cover={playlist.cover}/>
{/each}
</div>

{/if}

{#if $search_results.radios != null}
<h1>Radios</h1>

<div class="flex flex-row">
{#each $search_results.radios as radio}
    <CollectionCard type="Radio" view_type="radio" id={radio.id} name={radio.name} cover={"0"}/>
{/each}
</div>

{/if}

{#if $search_results.tracks == null && $search_results.radios == null && search_playlists == null && search_albums == null}
<h1>No search results.</h1>
{/if}

</div>

<style>

    h1 {
        @apply text-gray-300;
        @apply font-bold;
        @apply text-2xl;
    }

    table {
        @apply w-full;
        @apply m-3;
    }

    th {
        @apply font-thin;
        @apply text-gray-300;
    }

</style>