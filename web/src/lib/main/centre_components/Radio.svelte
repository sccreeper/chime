<script>

    import { active_view } from "../../stores";

    export let name;
    export let radio_id;

    function loadRadio() {
        active_view.set({name: "radio", id: radio_id})
    }

</script>

<p on:click={loadRadio} class="{($active_view.id == radio_id && $active_view.name == "radio") ? 'active' : ''}">{name}</p>

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