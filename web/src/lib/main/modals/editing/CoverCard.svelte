<script>

    import default_cover from "../../../../assets/no_cover.png";

    export let id = "";
    export let target_id = "";
    export let selected = "false";
    export let data_callback = (id) => {}

    async function clicked() {
        selected = "true";
        data_callback(id);

        await fetch("/api/edit/cover", {
            method: "POST",
            body: JSON.stringify({
                collection_id: target_id,
                cover_id: id,
            })
        })

    }

</script>

<img on:click={clicked} src="{id == "0" ? default_cover : `/api/collection/get_cover/${id}`}" data-selected="{selected}" alt="Album cover"/>

<style>

    img {
        @apply m-1;
        @apply border;
        @apply border-transparent;
        @apply w-24;
        @apply h-24;
        @apply transition-all;
        @apply cursor-pointer;
    }

    img:hover {
        @apply border-yellow-600;
    }

    img[data-selected=true] {
        @apply border-yellow-600;
    }

</style>