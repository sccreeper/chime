<script>
    import { page } from "$app/stores";
    import { CURRENT_COLLECTION_CONTEXT_KEY, PLAYER_CONTEXT_KEY } from "$lib/player";
    import { getContext, onMount } from "svelte";

    /** @type {import('$lib/player').ChimePlayer} */
    const player = getContext(PLAYER_CONTEXT_KEY)

    /** @type {import('$lib/player').ChimePlayer} */
    const {collectionId, playing, currentTrack} = getContext(PLAYER_CONTEXT_KEY)

    export let track_id = ""
    export let index = 0
    export let displayed_index = 0

    let mouse_over = false

    function click() {
        player.playCollection($page.data.collection, index, $page.data.collection_id)
    }

</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
<p on:mouseenter={() => mouse_over = true} on:mouseleave={() => mouse_over = false} on:click={click}>

    {#if !mouse_over && $page.data.collection_id == $collectionId && $playing && $currentTrack?.id == track_id}
    <i class="bi bi-volume-up-fill text-yellow-500"></i>
    {:else if !mouse_over}
    {displayed_index+1}
    {:else}
    <i class="bi bi-play-fill text-yellow-500"></i>
    {/if}

</p>

<style lang="postcss">

    * {
        cursor: pointer;
    }

    p {
        @apply text-gray-500;
        @apply h-9;
    }

    i {
        @apply text-xl;
    }

</style>