<script>
    import { get } from "svelte/store";
    import { active_view } from "../../stores";
    import no_cover from "../../../assets/no_cover.png";
    import { onMount } from "svelte";
    import { audio_source, playing } from "../../player";

    let cover_id = ""
    let name = ""
    let description = ""
    let url = ""
    let id = ""

    function load_details() {
        
        fetch(`/api/get_radio/${get(active_view).id}`, {
            method: "GET"
        }).then(response => response.json()).then(data => {

            cover_id = data.cover_id
            name = data.name
            description = data.description
            url = data.url
            id = get(active_view).id

        })

    }

    onMount(() => {
        load_details()
    })

    function handle_click() {
        // This radio is the radio being played.
        if (get(playing) && get(active_view).name == "radio" && get(active_view).id == id && get(audio_source).source == url) {
            playing.set(false)
        
        // This radio isn't being played, we can just set it normally.
        } else {
            audio_source.set({type: "radio", source: url})
        }
    }

    active_view.subscribe(load_details)

</script>

<div class="flex justify-center items-center w-full h-full">
    
    <div>
        <div class="flex flex-row items-center gap-3">
            <img src={cover_id == "0" || cover_id == "00" ? no_cover : `/api/collection/get_cover/${cover_id}`} width="200px" height="200px"/>
            <div class="flex flex-col gap-3">
                <h1 class="album-title">{name}</h1>
                <p class="text-gray-300">{description}</p>
            </div>
        </div>
        
        <button on:click={handle_click} class="mt-4">
            {#if $active_view.name == "radio" && $active_view.id == id && $playing && $audio_source.source == url}
            <i class="bi bi-pause-fill"></i> Pause
            {:else}
            <i class="bi bi-play-fill"></i> Play
            {/if}
        </button>
    
    </div>


</div>

<style>

.album-title {
    @apply text-gray-300;
    @apply font-bold;
    @apply text-5xl;
}

</style>

