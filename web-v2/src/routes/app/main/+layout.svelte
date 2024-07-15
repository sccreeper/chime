<script>
    import { browser } from '$app/environment';
    import { navigating } from '$app/stores';
    import { examineTrackId } from '$lib';
    import { getTrackMetadata } from '$lib/api/library';
    import BlankPage from '$lib/components/general/BlankPage.svelte';
    import ListDivider from '$lib/components/general/ListDivider.svelte';
    import ToggleButton from '$lib/components/general/ToggleButton.svelte';
    import LibraryItem from '$lib/components/main/library/LibraryItem.svelte';
    import { ChimePlayer, PLAYER_CONTEXT_KEY } from '$lib/player';
    import { getContext, setContext } from 'svelte';
    import defaultCover from '$lib/assets/no_cover.png';
    import HorizontalDivider from '$lib/components/general/HorizontalDivider.svelte';
    import { convertDuration, coverSizes } from '$lib/util';
    import MinorButton from '$lib/components/general/MinorButton.svelte';
    import { goto, invalidateAll } from '$app/navigation';

    /** @type {import('./$types').LayoutData} */
    export let data;

    /** @type {import('$lib/player').ChimePlayer} */
    setContext(PLAYER_CONTEXT_KEY, new ChimePlayer());

    /** @type {import('$lib/player').ChimePlayer} */
    const { playing, durationString, currentTimeString, duration, currentTime, shuffle, repeat, volume, currentTrack } = getContext(PLAYER_CONTEXT_KEY);
    /** @type {import('$lib/player').ChimePlayer} */
    const player = getContext(PLAYER_CONTEXT_KEY);

    /** @type {import('$lib/api/models').TrackMetadata} */
    let trackMetadata;

    /** @type {HTMLFormElement} */
    let searchForm;
    /** @type {number|undefined} */
    let searchScheduler = undefined;

    function scheduleSearch() {

        if (searchScheduler != undefined) {
            window.clearTimeout(searchScheduler)    
        }
        
        searchScheduler = window.setTimeout(() => {
            searchForm.requestSubmit()
        }, 100)

    }

    currentTrack.subscribe((val) => {
        if (val != undefined && browser) {
            examineTrackId.set(val?.id);   
        }
    })

    examineTrackId.subscribe(async (val) => {

        if (val == "") {
            return
        }

        trackMetadata = await getTrackMetadata(val)
    })

    // New radio and playlist

    function newPlaylist() {
        
        fetch("/api/collection/add", {
            method: "POST",
            body: JSON.stringify({
                "name": "Untitled Playlist",
                "description": "",
                "is_album": false,
            }),
        }).then((resp) => resp.json()).then(
            (resp) => {
                data.lib.playlists.push({
                    name: "Untitled Playlist",
                    cover_id: "0",
                    id: resp.id,
                })

                invalidateAll()
                goto(`/app/main/collection/${resp.id}`)
            }
        )

    }

    function newRadio() {
        
    }
</script>

<svelte:head>
    <title>Home - Chime</title>
</svelte:head>

{#if $navigating}

<div class="loader">
    <span>
    </span>
</div>

{/if}

<div class="flex flex-col h-full w-full">


    <!-- Search, navigation, settings etc. -->
    <div class="flex flex-row flex-nowrap items-center justify-center gap-3 h-20">

        <div class="w-1/5">        
            <form bind:this={searchForm} action="/app/main/search" autocomplete="off" data-sveltekit-keepfocus>
                <input class="w-full" type="text" name="query" placeholder="Search" on:input={scheduleSearch}/>
            </form>
        </div>

        <i class="bi bi-search text-gray-500"></i>

    </div>

    <!-- Central content -->

    <div class="flex h-full grow overflow-hidden">
        
        <div class="centre-left">

            <div class="m-2">

                <ListDivider text="Albums" icon="vinyl-fill"/>

                {#each data.lib.albums as item}
                <LibraryItem data={item} type="collection"/>
                {/each}

                <ListDivider text="Playlists" icon="list">
                    <MinorButton icon="plus" callback={newPlaylist}/> 
                </ListDivider>

                {#each data.lib.playlists as item}
                <LibraryItem data={item} type="collection"/>
                {/each}

                <ListDivider text="Radios" icon="broadcast"/>

                {#each data.lib.radios as item}
                <LibraryItem data={item} type="radio"/>
                {/each}

            </div>

        </div>

        <div class="centre">
            <slot />
        </div>

        <div class="centre-right">

            {#if $examineTrackId == "" || $examineTrackId == undefined || trackMetadata == undefined}
                <BlankPage icon="music-note-list" text="Not looking at anything right now"/>
            {:else}

                <div class="flex flex-col items-center text-center overflow-y-scroll overflow-x-hidden h-full w-full text-gray-500">

                    <img src={trackMetadata.cover_id == "0" ? defaultCover : `/api/collection/get_cover/${trackMetadata.cover_id}?width=${coverSizes.large}&height=${coverSizes.large}`} width="300" height="300" class="mt-3"/>

                    <h1 class="mt-2 text-2xl font-bold text-gray-300">{trackMetadata.title}</h1>
                    <p class="text-gray-300 mb-2">{trackMetadata.artist} <span class="text-yellow-600">●</span> {trackMetadata.album_name}</p>
                    <HorizontalDivider/>
                    <table class="text-left w-full text-sm m-2">
                        <colgroup>
                            <col span="1" style="width: 25%;">
                            <col span="1" style="width: 75%;">
                        </colgroup>

                        <tr><td class="font-semibold">Released:</td><td>{trackMetadata.released}</td></tr>
                        <tr><td class="font-semibold">Duration:</td><td>{convertDuration(trackMetadata.duration)}</td></tr>
                        <tr><td class="font-semibold">Format:</td><td>{trackMetadata.format}</td></tr>
                        <tr><td class="font-semibold">Original file:</td><td>{trackMetadata.original_file}</td></tr>
                        <tr><td class="font-semibold">File size:</td><td>{(trackMetadata.size / Math.pow(10, 6)).toFixed(2)} mb</td></tr>

                    </table>

                    <HorizontalDivider/>

                    <div class="w-full flex-col items-center gap-3 mt-3">
                        <MinorButton icon="download" callback={() => {}} hint="Download original file"/>
                        <MinorButton icon="plus-lg" callback={() => {}} hint="Add to collection"/>
                        <MinorButton icon="pencil-fill" callback={() => {}} hint="Edit"/>
                    </div>

                </div>

            {/if}

        </div>

    </div>

    <!-- Player -->

    <div>

        <div class="grid grid-rows-1 player-grid justify-items-center items-center flex-none m-2">

            <div class="grid grid-rows-2">
        
                <!-- Controls -->
        
                <div class="grid grid-rows-1 grid-cols-5 gap-3 items-center justify-items-center">
                    <ToggleButton callback={() => {}} icon="shuffle" bind:active={$shuffle}/>
                    <button class="control-button" on:click={()=>{ player.nextTrack("backward") }}><i class="bi bi-skip-backward"></i></button>
                    
                    <button 
                        class="control-button" 
                        on:click={() => { $playing = !$playing}}>
                            {#if !$playing}<i class="bi bi-play-fill"></i>{:else}<i class="bi bi-pause-fill"></i>{/if}
                    </button>
                    
                    <button class="control-button" on:click={() => { player.nextTrack("forward") }}><i class="bi bi-skip-forward"></i></button>
                    <ToggleButton callback={() => {}} icon="repeat" bind:active={$repeat}/>
                </div>
        
                <!-- Volume -->
        
                <div class="grid grid-rows-1 grid-cols-2 volume-grid items-center justify-items-center">
                    <i class="bi bi-{'volume-up'} mr-1 text-slate-400" on:click={()=>{}}></i>
                    <input class="seek" type="range" min="0" max="1.0" step="0.05" bind:value={$volume}/>
                </div>
        
        
            </div>
        
            <!-- Seek -->
        
            {#if false}
        
                <div class="flex justify-center items-center w-full h-full">
                    <p class="text-slate-400">Listening to internet radio.</p>
                </div>
        
            {:else}
        
            <div class="grid grid-cols-3 grid-rows-1 seek-grid w-full h-full items-center justify-items-center">
                <p class="duration select-none">{!browser ? '--:--' : $currentTimeString}</p>
                <input 
                    class="seek w-full" 
                    type="range" 
                    max={$duration}
                    bind:value={$currentTime}
                    on:change={
                        (e) => {
                            // @ts-ignore
                            player.seek(e.target.value)
                        }
                    }
                />
                <p class="duration select-none">{!browser ? '--:--' : $durationString}</p>
            </div>
            {/if}
        
        </div>

    </div>


</div>

<style lang="postcss">

    .centre-left {
        width: 15%;
        background-color: rgb(45, 53, 66);

        @apply h-full;
        @apply select-none;
        @apply overflow-y-scroll;

    }

    .centre-right {
        width: 20%;
        background-color: rgb(45, 53, 66);
    }

    .centre {
        width: 65%;
        background-color: rgb(38, 45, 56);
    }

    /* Loading bar */

    .loader {
        @apply w-full;
        @apply absolute;
        @apply bg-yellow-600;
        @apply h-1;
        @apply overflow-hidden;

        @apply left-0;
        @apply top-0;
    }

    @keyframes loader-animation {

        0% {
            @apply left-0;
            @apply right-full;
            @apply w-0;
        }

        10% {
            @apply left-0;
            @apply right-3/4;
            @apply w-1/4;
        }

        90% {
            @apply right-0;
            @apply left-3/4;
            @apply w-1/4;
        }

        100% {
            @apply left-full;
            @apply right-0;
            @apply w-0;
        }

    }

    .loader > span {
        @apply absolute;
        @apply h-full;
        @apply bg-yellow-500;
        animation: loader-animation 1s infinite;
    }

    /* Player CSS */

    .duration {
        @apply text-slate-400;
    }

    .player-grid {
        grid-template-columns: 2fr 8fr;
    }

    .volume-grid {
        grid-template-columns: 1fr 9fr;
    }

    .seek-grid {
        grid-template-columns: 1fr 4fr 1fr;
    }

    .seek {
        -webkit-appearance: none;
        background-color: transparent;
    }

    .seek::-webkit-slider-runnable-track {
        -webkit-appearance: none;
        appearance: none;

        background-color: transparent;
        border: none;
        box-shadow: none;

        @apply h-1;
        @apply bg-slate-500;
        @apply outline-none;
    }

    .seek::-moz-range-track {
        appearance: none;
        background-color: transparent;
        border: none;
        box-shadow: none;

        @apply h-1;
        @apply bg-slate-500;
        @apply outline-none;
    }


    .seek::-moz-range-thumb {
        -webkit-appearance: none;
        appearance: none;
        @apply h-3;
        @apply w-3;
        @apply bg-yellow-600;
        @apply outline-none;
        @apply border-none;
        cursor: pointer;
    }


    .seek::-webkit-slider-thumb {
        -webkit-appearance: none;
        @apply h-3;
        @apply w-3;
        @apply bg-yellow-600;
        @apply outline-none;
        @apply rounded-full;
        transform: translateY(-25%);
        cursor: pointer;
    }

    .control-button {

        border-radius: 50%;
        cursor: pointer;
        @apply w-10;
        @apply h-10;
        @apply text-lg;
    }

</style>