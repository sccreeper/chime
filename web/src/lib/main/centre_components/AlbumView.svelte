<script>
    import { active_view, collection_tracks } from "../../stores";
    import Track from "./Track.svelte";
    import no_cover_image from "../../../assets/no_cover.png";
    import HorizontalDivider from "../general/HorizontalDivider.svelte";
    import { closeModal, openModal } from "svelte-modals";
    import ConfirmModal from "../modals/ConfirmModal.svelte";
    import { get } from "svelte/store";
    import MinorButtonText from "../general/MinorButtonText.svelte";
    import { audio_source, playing, playing_collection, shuffle, viewing_tracks } from "../../player";
    import CollectionAdd from "../modals/CollectionAdd.svelte";
    import EditCollection from "../modals/editing/EditCollection.svelte";
    import { convertDurationLong } from "../../util";
    import EditCover from "../modals/editing/EditCover.svelte";

    let collection_title = "";
    let collection_cover_src = no_cover_image;

    let is_album = false;
    let collection_description = "";
    let actual_album_cover = no_cover_image;
    let is_protected = false;

    let title_font = "4.5vw";

    let cover_id = "";

    let collection_duration = "";

    function updateView(data) {
        if (data.name == "radio" || data.name == "search") {
            return
        } else {

            fetch(`/api/get_collection/${data.id}`, {
                method: "GET",
            })
                .then((response) => response.json())
                .then((data) => {
                    collection_title = data.title;
                    collection_tracks.set(data.tracks);


                    let track_ids = []

                    data.tracks.forEach(element => {
                        track_ids.push(element.id)
                    });

                    viewing_tracks.set(track_ids)

                    collection_description = data.description;
                    is_album = data.is_album
                    is_protected = data.protected

                    if (data.cover == "0" || data.cover == "00") {
                        collection_cover_src = ""
                    } else {
                        collection_cover_src = `/api/collection/get_cover/${data.cover}`
                    }

                    cover_id = data.cover;

                    if (collection_title.length > 35) {
                        title_font = "2vw"
                    } else if (collection_title.length > 25) {
                        title_font = "3vw"
                    } else if (collection_title.length > 15) {
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
        } else if ($collection_tracks.length == 0) {
            return
        } else if (get(shuffle)) {
            audio_source.set({type: "track", source: $collection_tracks[Math.floor(Math.random()*$collection_tracks.length)].id})

            playing_collection.set(get(active_view).id)
        } else {
            audio_source.set({type: "track", source: $collection_tracks[0].id})

            playing_collection.set(get(active_view).id)
        }

    }

    function addToCollection() {
        openModal(CollectionAdd, {id: get(active_view).id, type: "collection", exclude: get(active_view).id})
    }

    function editCollection() {
        openModal(EditCollection, {
            collection_id: get(active_view).id, 
            collection_description: collection_description, 
            collection_is_album: is_album, 
            collection_title: collection_title,
            data_callback: editCallback
        })
    }

    function editCallback(data) {
        
        collection_description = data.collection_description
        collection_title = data.collection_title

    }

    function editCover() {
        
        openModal(EditCover, {
            current_cover: cover_id,
            target_cover: get(active_view).id,
            data_callback: editCoverCallback,
            
        })

    }

    function editCoverCallback(id) {
        
        actual_album_cover = id == "0" ? no_cover_image : `/api/collection/get_cover/${id}?width=300&height=300`

    }

    $: actual_album_cover = collection_cover_src == "" ? no_cover_image : `${collection_cover_src}?width=300&height=300`

    $: {
        
        let duration = 0;
        $collection_tracks.map((e) => {duration += e.duration; console.log(duration)});

        collection_duration = convertDurationLong(duration);

    }

</script>

<div class="m-2 h-full overflow-y-scroll">
    <div class="bg-image" style={`background-image: url(${actual_album_cover});`}>
        <div class="album-title-container">
            <img on:click={() => {window.open(`${window.location.protocol}//${window.location.hostname}${actual_album_cover}`)}}
                src={actual_album_cover}
                class="album-cover"
            />

            <div class="flex flex-col gap-4 items-start">
                <h1 class="album-title" style="font-size: {title_font};">{collection_title}</h1>
                <p>{collection_description}</p>
                <p class="text-xs">{collection_duration}</p>
                <div class="flex flex-row items-center gap-3">
                    <button on:click={playCollection}>{#if $playing_collection == $active_view.id && $playing}<i class="bi bi-pause-fill"></i> Pause{:else}<i class="bi bi-play-fill"></i> Play{/if}</button> 
                    {#if !is_protected}
                        <MinorButtonText callback={deleteCollection} text="Delete" icon="trash-fill"/>
                        <MinorButtonText callback={editCollection} text="Edit" icon="pencil-fill"/>
                        <MinorButtonText callback={editCover} text="Edit cover" icon="image"/>
                    {/if}
                    {#if is_album}<MinorButtonText callback={addToCollection} text="Add to" icon="plus-lg"/>{/if}
                </div>
            </div>
        </div>
    </div>

    <div class="m-2"><HorizontalDivider/></div>

    <div>
        {#if $collection_tracks.length == 0}
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

            {#each $collection_tracks as track, i}
                <Track
                    index={i+1}
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
        @apply cursor-zoom-in;
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
    }

    .album-title-container {
        @apply flex flex-row items-center gap-4 backdrop-blur-lg backdrop-brightness-50 w-full;
    }

</style>
