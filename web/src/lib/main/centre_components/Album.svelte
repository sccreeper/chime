
<script>
    import { get } from "svelte/store";
    import { current_album } from "../../stores";


export let album_id;
export let name;

function open_album() {
    current_album.set(album_id)
}

</script>

<p on:click={open_album} class="{$current_album == album_id ? 'active' : ''}">{name}</p>

<style>

    p {
        @apply text-gray-300;
        @apply transition-all;
    }

    p:hover {
        @apply text-yellow-600;
    }

    .active {
        @apply text-yellow-600;
    }

</style>