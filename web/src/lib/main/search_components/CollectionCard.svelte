<script>
    import { active_view } from "../../stores";

    import no_cover from "../../../assets/no_cover.png"

    export var type = ""
    export var view_type = ""
    export var id = ""
    export var name = ""
    export var cover = ""

    function setView() {
        active_view.set({
            name: view_type,
            id: id,
        })
    }
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<div on:click={setView} title={name}>
    <img src={cover == "0" ? no_cover : `/api/collection/get_cover/${cover}`} width="108" height="108"/>
    <h1 class="">{name}</h1>
    <p>{type}</p>
</div>

<style>

    p {

        @apply text-xs;
        @apply font-semibold;
        @apply text-gray-300;

    }

    h1 {
        @apply overflow-hidden overflow-ellipsis whitespace-nowrap;
        @apply text-gray-300;
    }

    div { 
        @apply w-28;
        @apply h-40;
        @apply inline;
        background-color: rgb(45, 53, 66);
        @apply m-3;
        @apply p-2;
        @apply cursor-pointer;
        @apply transition-all;

        @apply border;
        border-color: rgb(45, 53, 66);
    }

    div:hover {

        @apply border;
        @apply border-yellow-600;

    }
</style>