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
<div class="" on:click={setView}>
    <img src={cover == "0" ? no_cover : `/api/collection/get_cover/${cover}`} width="64" height="64"/>
    <h1>{name}</h1>
    <p>{type}</p>
</div>

<style>

    div {
        @apply w-20; 
        @apply h-32;
        @apply rounded-md;
        @apply inline;
    }

</style>