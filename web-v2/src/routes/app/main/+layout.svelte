<script>
    import { browser } from '$app/environment';
    import { navigating } from '$app/stores';
    import BlankPage from '$lib/components/general/BlankPage.svelte';
    import ListDivider from '$lib/components/general/ListDivider.svelte';
    import ToggleButton from '$lib/components/general/ToggleButton.svelte';
    import LibraryItem from '$lib/components/main/library/LibraryItem.svelte';
    import { ChimePlayer, PLAYER_CONTEXT_KEY } from '$lib/player';
    import { getContext, setContext } from 'svelte';
    import { get, writable } from 'svelte/store';

    /** @type {import('./$types').LayoutData} */
    export let data;

    /** @type {import('$lib/player').ChimePlayer} */
    setContext(PLAYER_CONTEXT_KEY, new ChimePlayer());

    /** @type {import('$lib/player').ChimePlayer} */
    const { playing, durationString, currentTimeString, duration, currentTime } = getContext(PLAYER_CONTEXT_KEY);
    /** @type {import('$lib/player').ChimePlayer} */
    const player = getContext(PLAYER_CONTEXT_KEY);
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
    <div>
        <h1>Top Bar</h1>
    </div>

    <!-- Central content -->

    <div class="flex h-full grow overflow-hidden">
        
        <div class="centre-left">

            <div class="m-2">

                <ListDivider text="Albums" icon="vinyl-fill"/>

                {#each data.lib.albums as item}
                <LibraryItem name={item.name} item_id={item.id} type="collection"/>
                {/each}

                <ListDivider text="Playlists" icon="list"/>

                {#each data.lib.playlists as item}
                <LibraryItem name={item.name} item_id={item.id} type="collection"/>
                {/each}

                <ListDivider text="Radios" icon="broadcast"/>

                {#each data.lib.radios as item}
                <LibraryItem name={item.name} item_id={item.id} type="radio"/>
                {/each}

            </div>

        </div>

        <div class="grow centre">
            <slot />
        </div>

        <div class="centre-right">

            <BlankPage icon="music-note-list" text="Not looking at anything right now"/>

        </div>

    </div>

    <!-- Player -->

    <div>

        <div class="grid grid-rows-1 player-grid justify-items-center items-center flex-none m-2">

            <div class="grid grid-rows-2">
        
                <!-- Controls -->
        
                <div class="grid grid-rows-1 grid-cols-5 gap-3 items-center justify-items-center">
                    <ToggleButton callback={(x) => {}} icon="shuffle"/>
                    <button class="control-button" on:click={()=>{}}><i class="bi bi-skip-backward"></i></button>
                    
                    <button 
                        class="control-button" 
                        on:click={() => { $playing = !$playing}}>
                            {#if !$playing}<i class="bi bi-play-fill"></i>{:else}<i class="bi bi-pause-fill"></i>{/if}
                    </button>
                    
                    <button class="control-button" on:click={() => {}}><i class="bi bi-skip-forward"></i></button>
                    <ToggleButton callback={(x) => {}} icon="repeat"/>
                </div>
        
                <!-- Volume -->
        
                <div class="grid grid-rows-1 grid-cols-2 volume-grid items-center justify-items-center">
                    <i class="bi bi-{'volume-up'} mr-1 text-slate-400" on:click={()=>{}}></i><input class="seek" type="range" min="0" max="1.0" step="0.05"/>
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
        width: 450px;
        background-color: rgb(45, 53, 66);
    }

    .centre {
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