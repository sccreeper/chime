<script>
    import { track_metadata_view } from "../../stores";
    import { convertDuration } from "../../util";
    import TrackPlay from "./TrackPlay.svelte";

    export let id = "";
    export let title = "";
    export let artist = "";
    export let duration = 0;
    export let album_name = "";
    export let index = ""

    function changeView() {
        track_metadata_view.set(id)
    }
</script>

<tr on:click={changeView}>
    <td><TrackPlay index={index} track_id={id}/></td>
    <td class="font-semibold">{title}</td>
    <td class="text-xs">{artist}</td>
    <td class="text-xs">{album_name}</td>
    <td class="text-xs">{convertDuration(duration)}</td>
</tr>

<style>
    td {
        @apply cursor-pointer;
        @apply select-none;
        text-overflow: ellipsis;
        white-space: nowrap;
        @apply text-gray-500;
    }
</style>
