<script>
    import { active_view } from "../../stores";
    import Track from "./Track.svelte";
    import no_cover_image from "../../../assets/no_cover.png";
    import HorizontalDivider from "../general/HorizontalDivider.svelte";
    import { closeModal, openModal } from "svelte-modals";
    import ConfirmModal from "../modals/ConfirmModal.svelte";
    import { get } from "svelte/store";
    import MinorButtonText from "../general/MinorButtonText.svelte";
    import { audio_source, playing, playing_collection, shuffle, viewing_tracks } from "../../player";
    import CollectionAdd from "../modals/CollectionAdd.svelte";

    let album_title = "";
    let album_cover_src = no_cover_image;
    let tracks = [];
    let is_album = false;
    let album_description = "";
    let actual_album_cover = no_cover_image;

    let title_font = "4.5vw";

    function updateView(data) {
        if (data.name == "radio" || data.name == "search") {
            return
        } else {

            fetch(`/api/get_collection/${data.id}`, {
                method: "GET",
            })
                .then((response) => response.json())
                .then((data) => {
                    album_title = data.title;
                    tracks = data.tracks;


                    let track_ids = []

                    data.tracks.forEach(element => {
                        track_ids.push(element.id)
                    });

                    viewing_tracks.set(track_ids)

                    album_description = data.description;
                    is_album = (data.is_album == 1) ? true : false

                    if (data.cover == "0" || data.cover == "00") {
                        album_cover_src = ""
                    } else {
                        album_cover_src = `/api/collection/get_cover/${data.cover}`
                    }

                    if (album_title.length > 35) {
                        title_font = "2vw"
                    } else if (album_title.length > 25) {
                        title_font = "3vw"
                    } else if (album_title.length > 15) {
                        title_font = "3.5vw"
                    } else {
                        title_font = "4vw"
                    }

                });
        }
    }

    active_view.subscribe((value) => updateView(value));

    function deleteCollection() {
        
        openModal(ConfirmModal, {callback: (confirmed) => {

            if (confirmed) {
                
                fetch("/api/collection/delete", {
                    method: "POST",
                    body: JSON.stringify({collection_id: get(active_view).id})
                }).then((resp) => {
                    closeModal()
                })

            }

        }, message: "Are you sure you want to delete this collection. This will not delete it's contents."})

    }

    function playCollection() {

        if (get(playing) && get(playing_collection) == get(active_view).id) {
            playing.set(false)
        } else if (tracks.length == 0) {
            return
        } else if (get(shuffle)) {
            audio_source.set({type: "track", source: tracks[Math.floor(Math.random()*tracks.length)].id})

            playing_collection.set(get(active_view).id)
        } else {
            audio_source.set({type: "track", source: tracks[0].id})

            playing_collection.set(get(active_view).id)
        }

    }

    function addToCollection() {
        openModal(CollectionAdd, {id: get(active_view).id, type: "collection", exclude: get(active_view).id})
    }

    $: actual_album_cover = album_cover_src == "" ? no_cover_image : album_cover_src
</script>

<div class="m-2 h-full overflow-y-scroll">
    <div class="bg-image" style={`background-image: url(${actual_album_cover});`}>
        <div class="album-title-container">
            <img
                src={actual_album_cover}
                class="album-cover"
            />

            <div class="flex flex-col gap-4 items-start">
                <h1 class="album-title" style="font-size: {title_font};">{album_title}</h1>
                <p>{album_description}</p>
                <div class="flex flex-row items-center gap-3">
                    <button on:click={playCollection}>{#if $playing_collection == $active_view.id && $playing}<i class="bi bi-pause-fill"></i> Pause{:else}<i class="bi bi-play-fill"></i> Play{/if}</button> 
                    {#if $active_view.id != "1"}<MinorButtonText callback={deleteCollection} text="Delete" icon="trash-fill"/>{/if}
                    {#if is_album}<MinorButtonText callback={addToCollection} text="Add to" icon="plus-lg"/>{/if}
                </div>
            </div>
        </div>
    </div>

    <div class="m-2"><HorizontalDivider/></div>

    <div>
        {#if tracks.length == 0}
            <p>
                No tracks in {is_album ? "album" : "playlist"}. Add from the
                lefthand sidebar or upload files.
            </p>
        {:else}
            <table>
                <colgroup>
                    <col span="1" style="width: 5%;">
                    <col span="1" style="width: 35%;">
                    <col span="1" style="width: 15%;">
                    <col span="1" style="width: 35%;">
                    <col span="1" style="width: 10%;">
                 </colgroup>

                <tr class="text-left">
                    <th>No.</th>
                    <th>Title</th>
                    <th>Artist</th>
                    <th>Album</th>
                    <th>Duration</th>
                </tr>

            {#each tracks as track, i}
                <Track
                    index={(i+1).toString()}
                    id={track.id}
                    title={track.name}
                    artist={track.artist}
                    duration={track.duration}
                    album_name={track.album_name}
                />
            {/each}

            </table>
        {/if}
    </div>
</div>

<style>
    .album-title {
        @apply text-gray-300;
        @apply font-bold;
    }

    .album-cover {
        width: 256px;
        min-width: 128px;
    }

    table {
        width: 100%;
    }

    th {
        @apply font-thin;
    }

    .bg-image {
        @apply w-full;
        @apply bg-no-repeat;
        @apply bg-cover;    
        @apply bg-center;
        @apply -z-20;
    }

    .album-title-container {
        @apply flex flex-row items-center gap-4 backdrop-blur-lg backdrop-brightness-50 w-full -z-50;
    }

</style>
