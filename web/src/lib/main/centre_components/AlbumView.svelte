<script>
    import { current_album, session_object } from "../../stores";
    import { get } from "svelte/store";
    import Track from "./Track.svelte";
    import no_cover_image from "../../../assets/no_cover.png";

    let album_title = "";
    let album_cover_src = "";
    let tracks = [];
    let is_album = false;
    let album_description = "";

    let title_font = "4.5vw";

    function updateView(album_id) {
        fetch(`/api/get_collection/${album_id}`, {
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

                if (album_title.length > 100) {
                    title_font = "1vw";
                } else if (album_title.length > 50) {
                    title_font = "3vw";
                } else if (album_title.length <= 25) {
                    title_font = "4.5vw";
                }

            });
    }

    

    current_album.subscribe((value) => updateView(value));
</script>

<div class="m-2">
    <div class="flex flex-row items-center gap-4">
        <img
            src={album_cover_src == "" ? no_cover_image : album_cover_src}
            class="album-cover"
        />

        <div class="flex flex-col gap-4 items-end">
            <h1 class="album-title" style="font-size: {title_font};">{album_title}</h1>
            <p>{album_description}</p>
        </div>
    </div>

    <hr class="m-3 h-px border-none bg-gray-600" />

    <div>
        {#if tracks.length == 0}
            <p>
                No tracks in {is_album ? "album" : "playlist"}. Add from the
                lefthand sidebar or upload files.
            </p>
        {:else}
            <table>
                <tr>
                    <th>Title</th>
                    <th>Artist</th>
                    <th>Album</th>
                    <th>Duration</th>
                </tr>
            </table>

            {#each tracks as track}
                <Track
                    id={track.id}
                    title={track.name}
                    artist={track.artist}
                    duration={track.duration}
                    album_name={track.album_name}
                />
            {/each}
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
</style>
