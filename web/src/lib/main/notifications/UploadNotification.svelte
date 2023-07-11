<script>
    import { onMount } from "svelte";
    import { removeNotification } from "../../notifications";


    export let progress = 0
    export let count = 0
    export let finished = false;

    let element;

    // Required attribute
    export let id = 0;
    
    if (finished) {
        console.log("lol 2")
        setTimeout(() => {
            removeNotification(id)
        }, 5000)
    }
    
    $: () => {
        console.log(progress)
        element.style.background = `linear-gradient(to right, rgb(235, 179, 8), rgb(235, 179, 8) ${Math.floor((progress / count) * 100)}%, rgb(202, 138, 4) ${Math.floor((progress / count) * 100)}%, rgb(202, 138, 4));`
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

    <div bind:this={element} class="w-full h-1 bg-yellow-600"></div>
</div>