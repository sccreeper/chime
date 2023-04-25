<script>
    import { get } from "svelte/store";
    import HorizontalDivider from "../general/HorizontalDivider.svelte";
    import { track_metadata_view } from "../../stores";
    import MinorButton from "../general/MinorButton.svelte";
    import { convertDuration } from "../../util";
    import { openModal } from "svelte-modals";
    import CollectionAdd from "../modals/CollectionAdd.svelte";

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

    } 

    // Load metdata when track details change.

    function loadMetadata() {

        fetch(`/api/get_track_metadata/${get(track_metadata_view)}`, {
            method : "GET"
        }).then(response => response.json()).then(data => {
            metadata = data
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
        openModal(CollectionAdd, {track_id: get(track_metadata_view)})
    }

</script>

{#if $track_metadata_view == null || $track_metadata_view == ""}

    <p>Not listening to or examining any tracks.</p>

{:else}

    <div class="flex flex-col items-center text-center m-2 overflow-y-scroll h-full">
        <img src={`/api/collection/get_cover/${metadata.cover_id}`} width="300" height="300"/>
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
            <MinorButton icon="download" callback={download}/>
            <MinorButton icon="plus-lg" callback={addToPlaylist}/>
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