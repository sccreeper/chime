<script>
    import BlankPage from "$lib/components/general/BlankPage.svelte";
import HorizontalDivider from "$lib/components/general/HorizontalDivider.svelte";
import ListDivider from "$lib/components/general/ListDivider.svelte";
import SearchCollection from "$lib/components/main/search/SearchCollection.svelte";
    import SearchTrack from "$lib/components/main/search/SearchTrack.svelte";


    /**
     * @type {import("./$types").PageData}
     */
    export let data;

    /** @type {import("$lib/api/models").SearchCollection[]} */
    let searchAlbums = [];
    /** @type {import("$lib/api/models").SearchCollection[]} */
    let searchPlaylists = [];

    $: {

        searchAlbums = []
        searchPlaylists = []

        data.searchResults.collections.forEach(element => {
            if (element.is_album) {
                searchAlbums.push(element)
            } else {
                searchPlaylists.push(element)
            }
        });

    }

</script>

<div class="overflow-y-scroll h-full">

{#if searchAlbums.length == 0 && searchPlaylists.length == 0 && data.searchResults.tracks.length == 0}

    {#if data.query == undefined}
        <BlankPage text="Enter a search query" icon="search"/>
    {:else}
        <BlankPage text={`No results for query "${data.query}"`} icon="search"/>
    {/if}

{/if}

{#if searchAlbums.length != 0}

    <div class="search-container">
    
        <ListDivider icon="vinyl-fill" text="Albums"/>

        <div class="search-results">

            {#each searchAlbums as c}
            <SearchCollection data={c}/>
            {/each}

        </div>

    </div>

    <HorizontalDivider/>

{/if}

{#if searchPlaylists.length != 0}

    <div class="search-container">
    
    <ListDivider icon="music-note-list" text="Playlists"/>

        <div class="search-results">

            {#each searchPlaylists as c}
            <SearchCollection data={c}/>
            {/each}

        </div>

    </div>

    <HorizontalDivider/>

{/if}

{#if data.searchResults.tracks.length != 0}

    <div class="search-container">
    
    <ListDivider icon="music-note" text="Tracks"/>

    <table>

        <colgroup>
        <col span="1" style="width: 3%;"> <!-- Cover -->
        <col span="1" style="width: 50%;"> <!-- Title -->
        <col span="1" style="width: 27%;"> <!-- Artist -->
        <col span="1" style="width: 20%;"> <!-- Duration -->
        </colgroup>

        <tr class="text-left">
            <th><i class="bi bi-image text-gray-500"></i></th>
            <th>Title</th>
            <th>Artist</th>
            <th>Duration</th>
        </tr>

        {#each data.searchResults.tracks as t }
            
            <SearchTrack data={t}/>

        {/each}

    </table>

    </div>

    <HorizontalDivider/>

{/if}

</div>

<style lang="postcss">

    .search-results {

        @apply flex flex-row gap-3 flex-nowrap;
        @apply overflow-x-scroll;
        @apply overflow-y-hidden;
        @apply w-full;

    }

    .search-container { 
        @apply m-3;
    }

    th {
        @apply font-thin;
    }

</style>