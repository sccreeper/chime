<script>
    import { get } from "svelte/store";
    import { album_list, session_object } from "../../stores";
    import { onMount } from "svelte";
    import Album from "./Album.svelte";
    import Radio from "./Radio.svelte";

    function loadAlbums() {

        fetch("/api/get_collections", {
            method: "GET",
        }).then((response) => response.json()).then((data) => {
            album_list.set(data)
        })

    }

    onMount(() => {
        loadAlbums()
    })

</script>

<div class="album-list">

<strong><i class="bi bi-vinyl-fill"></i> Albums</strong>

{#if $album_list.albums.length == 0}
<p class="text-gray-300">No Albums</p>
{:else}
{#each $album_list.albums as album}

<Album album_id={album.id} name={album.name}/>

{/each}
{/if}

<strong><i class="bi bi-list"></i> Playlists</strong>

{#if $album_list.playlists.length == 0}
<p class="text-gray-300">No playlists</p>
{:else}

<ul>

{#each $album_list.playlists as playlist}

<Album album_id={playlist.id} name={playlist.name}/>

{/each}

</ul>

{/if}
<strong><i class="bi bi-broadcast"></i> Radio</strong>

{#if $album_list.radios.length == 0}
<p class="text-gray-300">No radios</p>
{:else}
<ul>

{#each $album_list.radios as radio}

<Radio radio_id={radio.id} name={radio.name}/>

{/each}

</ul>
{/if}


</div>

<style>

.album-list {
    @apply h-full;
    @apply select-none;
    @apply m-2;
    @apply overflow-y-scroll;

}

strong {

    @apply text-yellow-500;

}

</style>