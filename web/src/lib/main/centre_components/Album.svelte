<script>
import { active_view } from "../../stores";
    import { goto } from "../../util";


export let album_id;
export let name;

function open_album() {
    active_view.set({name: "album", id: album_id})
    goto(`/collection/${album_id}`, "album", album_id)
}

</script>

<p on:click={open_album} class="{($active_view.id == album_id && $active_view.name == "album") ? 'active' : ''}">{name}</p>

<style>

    p {
        @apply text-gray-300;
        @apply transition-all;
        cursor: pointer;
        text-overflow: ellipsis;
    }

    p:hover {
        @apply text-yellow-600;
    }

    .active {
        @apply text-yellow-600;
    }

</style>