<script>
    import { active_view } from "../../stores";
    import Track from "./Track.svelte";
    import no_cover_image from "../../../assets/no_cover.png";
    import HorizontalDivider from "../general/HorizontalDivider.svelte";

    let album_title = "";
    let album_cover_src = "";
    let tracks = [];
    let is_album = false;
    let album_description = "";

    let title_font = "4.5vw";

    function updateView(data) {
        fetch(`/api/get_collection/${data.id}`, {
            method: "GET",
        })
            .then((response) => response.json())
            .then((data) => {
                album_title = data.title;
                tracks = data.tracks;
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

    

    active_view.subscribe((value) => updateView(value));
</script>

<div class="m-2">
    <div class="flex flex-row items-center gap-4">
        <img
            src={album_cover_src == "" ? no_cover_image : album_cover_src}
            class="album-cover"
        />

        <div class="flex flex-col gap-4 items-start">
            <h1 class="album-title" style="font-size: {title_font};">{album_title}</h1>
            <p>{album_description}</p>
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

</style>
