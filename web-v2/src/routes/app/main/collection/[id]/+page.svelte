<script>

    import no_cover from "$lib/assets/no_cover.png";
    import MinorButtonText from "$lib/components/general/MinorButtonText.svelte";
    import { convertDurationLong } from "$lib/util";

    /**
     * @type {import('./$types').PageData}
     */
    export let data

    $: actual_album_cover = data.collection.cover == "0" ? no_cover : `/api/collection/get_cover/${data.collection.cover}?width=300&height=300`
    
    /** @type {string} */
    let total_duration;
    $: {

        let duration = 0;
        data.collection.tracks.map((e) => {duration += e.duration})

        total_duration = convertDurationLong(duration)

    }

    // Play

    function playCollection() {
        
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

                    <button on:click={playCollection}><i class="bi bi-play-fill">Play</i></button>
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

</style>