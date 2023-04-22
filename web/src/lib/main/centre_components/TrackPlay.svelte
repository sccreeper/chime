<script>
    import { track_id } from "../../player";


    export let _track_id = ""
    export let index = ""

    let mouse_over = false

    function click() {
        track_id.set(_track_id)
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