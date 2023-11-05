<script>
    import { get } from "svelte/store";
    import HorizontalDivider from "../general/HorizontalDivider.svelte";
    import { active_view, collection_tracks, track_metadata_view } from "../../stores";
    import MinorButton from "../general/MinorButton.svelte";
    import { convertDuration } from "../../util";
    import { openModal } from "svelte-modals";
    import CollectionAdd from "../modals/CollectionAdd.svelte";
    import BlankPage from "../general/BlankPage.svelte";
    import EditTrack from "../modals/editing/EditTrack.svelte";
    import default_cover from "../../../assets/no_cover.png";

    // Metadata object

    let metadata = {

        title: "",
        album_name: "",
        cover_id: "",
        artist: "",
        original_file: "",
        format: "",
        duration: 0,
        released: 0,
        size: 0,
        album_id: "",
        previous_album_id: ""

    } 

    // Load metdata when track details change.

    function loadMetadata() {

        fetch(`/api/get_track_metadata/${get(track_metadata_view)}`, {
            method : "GET"
        }).then(response => response.json()).then(data => {
            metadata = data
            metadata.previous_album_id = metadata.album_id
        })

    }

    track_metadata_view.subscribe(loadMetadata)

    // Callback for downloading, deleting, editing etc.

    function download() {
        
        fetch(`/api/download_original/${get(track_metadata_view)}`, {
            method: "GET"
        }).then(response => response.arrayBuffer()).then(data => {

            let blob = new Blob([data])

            let link = document.createElement("a")
            link.href = URL.createObjectURL(blob)
            link.download = metadata.original_file

            link.click();

        })

    }

    function addToPlaylist() {
        openModal(CollectionAdd, {id: get(track_metadata_view), type: "track", exclude: ""})
    }

    function edit() {
        openModal(EditTrack, {
            track_title: metadata.title, 
            track_artist: metadata.artist, 
            track_released: metadata.released.toString(), 
            track_album_id: metadata.album_id,
            data_callback: (data) => {
                // Update data in UI
                metadata.title = data.title
                metadata.artist = data.artist
                metadata.artist = data.released
                metadata.album_id = data.album_id
                metadata.album_name = data.album_name

                // Update data in album view
                if (metadata.previous_album_id == metadata.album_id 
                && get(active_view).id == metadata.album_id 
                && get(active_view).name == "album") {
                    
                    let new_tracks = []

                    $collection_tracks.forEach((e) => {
                        let t = e

                        if (t.id == $track_metadata_view) {
                            t.name = data.title
                        }

                        new_tracks.push(t)
                    })

                    $collection_tracks = new_tracks


                }

            }
        })
    }

</script>

{#if $track_metadata_view == null || $track_metadata_view == ""}

<BlankPage icon="music-note-list" text="Not examining any tracks"/>

{:else}

    <div class="flex flex-col items-center text-center m-2 overflow-y-scroll h-full">
        <img src={metadata.cover_id == "0" ? default_cover : `/api/collection/get_cover/${metadata.cover_id}?width=300&height=300`} width="300" height="300"/>
        <h1 class="mt-2">{metadata.title}</h1>
        <p class="text-gray-300">{metadata.artist} <span class="dot">●</span> {metadata.album_name}</p>
        <HorizontalDivider/>
        <table class="text-left w-full">
            <colgroup>
                <col span="1" style="width: 25%;">
                <col span="1" style="width: 75%;">
             </colgroup>

            <tr><td class="header">Released:</td><td>{metadata.released}</td></tr>
            <tr><td class="header">Duration:</td><td>{convertDuration(metadata.duration)}</td></tr>
            <tr><td class="header">Format:</td><td>{metadata.format}</td></tr>
            <tr><td class="header">Original file:</td><td>{metadata.original_file}</td></tr>
            <tr><td class="header">File size:</td><td>{(metadata.size / Math.pow(10, 6)).toFixed(2)} mb</td></tr>

        </table>

        <HorizontalDivider/>

        <div class="w-full flex-col items-center">
            <MinorButton icon="download" callback={download} hint="Download original file"/>
            <MinorButton icon="plus-lg" callback={addToPlaylist} hint="Add to collection"/>
            <MinorButton icon="pencil-fill" callback={edit} hint="Edit"/>
        </div>

    </div>

{/if}

<style>

    * {
        @apply text-gray-500;
    }

    h1 {
        @apply text-2xl;
        @apply font-bold;
        @apply text-gray-300;
    }

    table {
        @apply text-sm;
    }

    .header {
        @apply font-semibold;
    }

    .dot {
        @apply text-yellow-600;
    }

</style>