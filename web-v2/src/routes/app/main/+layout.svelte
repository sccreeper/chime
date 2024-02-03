<script>
    import { navigating } from '$app/stores';
    import BlankPage from '$lib/components/general/BlankPage.svelte';
    import ListDivider from '$lib/components/general/ListDivider.svelte';
    import LibraryItem from '$lib/components/main/library/LibraryItem.svelte';

    /** @type {import('./$types').LayoutData} */
    export let data;
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

            <BlankPage icon="music-note-list" text="Not examining anything right now"/>

        </div>

    </div>

    <!-- Player -->

    <div>

        <div class="h-28">
            <p>Player</p>
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

</style>