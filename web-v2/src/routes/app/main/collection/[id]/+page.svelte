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
    import { afterNavigate, invalidateAll } from "$app/navigation";
    import Dialog from "$lib/components/dialogs/Dialog.svelte";
    import MinorButton from "$lib/components/general/MinorButton.svelte";
    import { applyAction } from "$app/forms";

    /**
     * @type {import('./$types').PageData}
     */
    export let data

    /** @type {import('$lib/player').ChimePlayer} */
    const player = getContext(PLAYER_CONTEXT_KEY)

    /** @type {import('$lib/player').ChimePlayer} */
    const {playing, collectionId} = getContext(PLAYER_CONTEXT_KEY)

    let beingEdited = false;
    let collectionTitle = data.collection.title;
    let collectionDescription = data.collection.description;
    afterNavigate(() => {
        beingEdited = false
        collectionTitle = data.collection.title;
        collectionDescription = data.collection.description;
    })

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
    /** @type {{index: number, track: import('$lib/api/models').CollectionTrack}[][]} */
    let discs;
    $: {

        discs = [];

        if (!data.collection.is_album) {
            
            discs.push([])
            discs[0].push(...data.collection.tracks.map((val, index) => ({track: val, index: index})))

        } else {

            // Edge cases

            if (data.collection.tracks.length == 0) {
                discs = [];
            }

            let disc_count = 0;

            discs.push([]);

            // Figure out how many discs there are in total
            // This is required because the data returned from the server isn't in order
            for (let i = 0; i < data.collection.tracks.length; i++) {
                const element = data.collection.tracks[i];

                if (((element.disc == 0 ? 1 : element.disc)-1) > disc_count) {
                    disc_count++
                    discs.push([])
                }   
            }

            for (let i = 0; i < data.collection.tracks.length; i++) {
                const element = data.collection.tracks[i];

                discs[(element.disc == 0 ? 1 : element.disc)-1].push({index: i, track: element})
                
            }

        }

        discs = discs;


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

    let showAddToCollectionModal = false;
    /**
     * @param {string} id
     */
    async function addToCollection(id) {

        await fetch("/api/collection/add_collection", {
            method: "POST",

            body: JSON.stringify(
                {
                    "source" : data.collection_id,
                    "destination" : id
                }
            )
        })

        showAddToCollectionModal = false;

    }

    async function editCollection() {

        if (beingEdited) {
            
            // Apply changes

            await fetch("/api/edit/collection", {
                method: "POST",
                body: JSON.stringify({
                    collection_id: data.collection_id,
                    name: collectionTitle,
                    description: collectionDescription,
                    is_album: data.collection.is_album,
                })
            })

            invalidateAll();
            beingEdited = false;

        } else {

            beingEdited = true;

        }

    }

    function editCover() {
        
    }

    function deleteCollection() {
        
    }

</script>

<svelte:head>
    <title>{data.collection.title}{data.collection.is_album && !data.collection.protected ? ` - ${data.collection.tracks[0].artist}` : ''} - Chime</title>
</svelte:head>

<!-- Dialogs -->

<Dialog bind:showModal={showAddToCollectionModal} title="Add to collection" icon="plus-lg">

    <div class="body-div" slot="body">

        {#each data.lib.playlists as playlist }
            <p on:click={() => {addToCollection(playlist.id)}}>{playlist.name}</p>
        {/each}

    </div>

</Dialog>

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

                {#if beingEdited}
                    
                    <input class="album-title" type="text" bind:value={collectionTitle} placeholder="Title">
                    <input class="album-description" type="text" bind:value={collectionDescription} placeholder="Description">

                {:else}
                
                <h1 class="album-title">{data.collection.title}</h1>
                <p class="album-description">{data.collection.description}</p>

                {/if}

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
                        <MinorButtonText callback={editCollection} text={beingEdited ? "Apply" : "Edit"} icon={beingEdited ? "check-lg" : "pencil-fill"}/>
                        <MinorButtonText callback={editCover} text="Edit cover" icon="image"/>
    
                    {/if}
    
                    {#if data.collection.is_album}<MinorButtonText callback={() => {showAddToCollectionModal = true}} text="Add to" icon="plus-lg"/>{/if}

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
                <col span="1" style="width: 35%;"> <!-- album -->
                <col span="1" style="width: 10%;"> <!-- favourite+duration -->
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

                {#each disc as track, j }
                    
                    <Track track={track.track} index={track.index} displayed_index={j}/>

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

    .body-div p {
        @apply text-slate-400;
        @apply transition-all;
        @apply cursor-pointer;
        @apply select-none;
    }

    .body-div p:hover {
        @apply text-yellow-500;
    }

</style>