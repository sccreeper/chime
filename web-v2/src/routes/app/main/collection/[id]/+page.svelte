<script>

    import no_cover from "$lib/assets/no_cover.png";
    import MinorButtonText from "$lib/components/general/MinorButtonText.svelte";
    import HorizontalDivider from "$lib/components/general/HorizontalDivider.svelte";
    import { convertDurationLong, coverSizes } from "$lib/util";
    import BlankPage from "$lib/components/general/BlankPage.svelte";
    import Disc from "$lib/components/main/collection/Disc.svelte";
    import Track from "$lib/components/main/collection/Track.svelte";
    import { getContext } from "svelte";
    import { PLAYER_CONTEXT_KEY } from "$lib/player";

    /**
     * @type {import('./$types').PageData}
     */
    export let data

    /** @type {import('$lib/player').ChimePlayer} */
    const player = getContext(PLAYER_CONTEXT_KEY)

    /** @type {import('$lib/player').ChimePlayer} */
    const {playing, collectionId} = getContext(PLAYER_CONTEXT_KEY)

    // Reactive variables in UI.

    $: actual_album_cover = data.collection.cover == "0" ? no_cover : `/api/collection/get_cover/${data.collection.cover}?width=${coverSizes.medium}&height=${coverSizes.medium}`
    
    // Find out total duration of collection in hh:mm
    /** @type {string} */
    let total_duration;
    $: {

        let duration = 0;
        data.collection.tracks.map((e) => {duration += e.duration})

        total_duration = convertDurationLong(duration)

    }

    // Sort tracks into 2D array, each 1D array being tracks on a different disc.
    /** @type {{index: number, track: import('$lib/api/api').CollectionTrack}[][]} */
    let discs;
    $: {

        discs = [];

        if (!data.collection.is_album) {
            
            discs.push([])
            discs[0].push(...data.collection.tracks.map((val, index) => ({track: val, index: index})))

        } else {

            let disc_count = 0;

            discs.push([]);

            for (let i = 0; i < data.collection.tracks.length; i++) {
                const element = data.collection.tracks[i];

                if (element.disc-1 > disc_count) {
                    
                    discs.push([])
                    disc_count++

                    discs[disc_count].push({index: i, track: element})

                    continue

                }

                discs[disc_count].push({index: i, track: element})
                
            }

        }

    }

    // Component methods.

    // Play

    function playCollection() {

        if ($playing && data.collection_id == $collectionId) {
            
            $playing = !$playing

        } else {

            if (data.collection_id == $collectionId) {
                $playing = !$playing
            } else {
                player.playCollection(data.collection, 0, data.collection_id)
            }

        }
    
    }

    // Edit functions

    function addToCollection() {
        
    }

    function editCollection() {
        
    }

    function editCover() {
        
    }

    function deleteCollection() {
        
    }

</script>

<svelte:head>
    <title>{data.collection.title}{data.collection.is_album && !data.collection.protected ? ` - ${data.collection.tracks[0].artist}` : ''} - Chime</title>
</svelte:head>

<div class="h-full overflow-y-scroll">

    <div class="bg-image" style={`background-image: url("${actual_album_cover}")`}>

        <div class="album-title-container">

            <!-- Cover Image -->

            <img on:click={() => {window.open(`${window.location.protocol}//${window.location.port == '' ? ":" + window.location.port : ''}${actual_album_cover}`)}}
            src={actual_album_cover}
            class="album-cover" alt="Cover for {data.collection.title}"
            />

            <!-- Text and buttons-->

            <div class="flex flex-col gap-4 items-start">

                <!-- Title, description, duration -->

                <h1 class="album-title">{data.collection.title}</h1>
                <p>{data.collection.description}</p>
                <p class="text-xs">{total_duration}</p>

                <!-- Buttons -->

                <div class="flex flex-row items-center gap-3">

                    {#if data.collection_id == $collectionId && $playing }
                    <button on:click={playCollection}><i class="bi bi-pause-fill"></i> Pause</button>
                    {:else}
                    <button on:click={playCollection}><i class="bi bi-play-fill"></i> Play</button>
                    {/if}
                    
                    {#if !data.collection.protected}
    
                        <MinorButtonText callback={deleteCollection} text="Delete" icon="trash-fill"/>
                        <MinorButtonText callback={editCollection} text="Edit" icon="pencil-fill"/>
                        <MinorButtonText callback={editCover} text="Edit cover" icon="image"/>
    
                    {/if}
    
                    {#if data.collection.is_album}<MinorButtonText callback={addToCollection} text="Add to" icon="plus-lg"/>{/if}

                </div>

            </div>

        </div>

    </div>

    <HorizontalDivider/>

    <!-- Collection table -->

    <div>

        {#if data.collection.tracks.length == 0}
            
            <BlankPage icon="folder2-open" text="Collection is empty"/>

        {:else}

        <table>

            <colgroup>
                <col span="1" style="width: 5%;"> <!-- no. -->
                <col span="1" style="width: 35%;"> <!-- title -->
                <col span="1" style="width: 15%;"> <!-- artist -->
                <col span="1" style="width: 35%;"> <!-- duration -->
                <col span="1" style="width: 10%;"> <!-- favourite -->
            </colgroup>

            <tr class="text-left">
                <th>No.</th>
                <th>Title</th>
                <th>Artist</th>
                <th>Album</th>
                <th>Duration</th>
            </tr>

            {#each discs as disc, i}
                
                {#if data.collection.is_album}
                    
                    <Disc number={i}/>

                {/if}

                {#each disc as track }
                    
                    <Track track={track.track} index={track.index}/>

                {/each}

            {/each}

        </table>

        {/if}

    </div>

</div>

<style lang="postcss">

    .bg-image {
        @apply w-full;
        @apply bg-no-repeat;
        @apply bg-cover;    
        @apply bg-center;
    }

    .album-title-container {
        @apply flex flex-row items-center gap-4 backdrop-blur-lg backdrop-brightness-50 w-full;
    }

    .album-cover {
        width: 256px;
        min-width: 128px;
        @apply cursor-zoom-in;
    }

    .album-title {

        @apply whitespace-nowrap overflow-hidden;
        min-width: 0;
        max-width: 100%;

        @apply text-gray-300;
        @apply font-bold;
        font-size: clamp(1rem, 2vw, 2rem);
    }

    table {
        max-width: 100%;
        @apply overflow-hidden;
        @apply m-2;
    }

    th {
        @apply font-thin;
    }

</style>