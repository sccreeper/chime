<script>
    import { get } from "svelte/store";
    import { audio_source, playing_collection } from "../../player";
    import { active_view } from "../../stores";


    export let track_id = ""
    export let index = 0

    let mouse_over = false

    function click() {
        playing_collection.set(get(active_view).id)

        audio_source.set({type: "track", source: track_id})
    }

</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<p on:mouseenter={() => mouse_over = true} on:mouseleave={() => mouse_over = false} on:click={click}>

    {#if !mouse_over}
    {index}
    {:else}
    <i class="bi bi-play-fill text-yellow-500"></i>
    {/if}

</p>

<style>

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