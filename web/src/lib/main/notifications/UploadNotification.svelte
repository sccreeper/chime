<script>
    import { removeNotification } from "../../notifications";


    export let progress = 0
    export let count = 0
    export let finished = false;

    // Required attribute
    export let id = 0;
    
    $: () => {
        if (finished) {
        console.log("lol 2")
        setTimeout(() => {
            removeNotification(id)
        }, 5000)
    }
    }

</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<div class="bg-gray-800 border border-yellow-600 w-full h-11 cursor-pointer" on:click={() => {removeNotification(id)}}>
    <div class="flex flex-row items-center w-full gap-2">

        <i class="bi bi-upload"/>

        {#if finished}
            <h1>Finished uploading {count} files</h1>
        {:else}
            <h1>Uploading {progress}/{count}</h1> 
        {/if}
    
    </div>

    <div style="--progress: {`${Math.floor((progress / count) * 100)}%`}" class="progress-bar"></div>
</div>

<style>
    .progress-bar {
        @apply w-full bg-yellow-600 h-2;
    }

    .progress-bar:after {
        content: "";
        @apply top-0 left-0 bottom-0 right-0;
        @apply block;
        width: var(--progress);
        @apply h-2 bg-yellow-400;
    }
</style>