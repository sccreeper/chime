<script>
    import { onMount } from "svelte";
    import { removeNotification } from "../../notifications";


    export let text = ""
    export let icon = ""
    export let expiry = 0

    // Required attribute
    export let id = 0

    onMount(() => {
        setTimeout(() => {
            removeNotification(id)
        }, expiry)
    })

</script>


<!-- svelte-ignore a11y-click-events-have-key-events -->
<div class="flex flex-row items-center bg-gray-800 border border-yellow-500 w-full h-11 gap-2 cursor-pointer" on:click={() => {removeNotification(id)}}>
    <i class="bi bi-{icon}"/>
    <h1 class="">{text}</h1>
</div>