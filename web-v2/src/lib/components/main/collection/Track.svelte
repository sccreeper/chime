<script>
    import { convertDuration } from "$lib/util";


    import FavouriteButton from "./FavouriteButton.svelte";
    import TrackPlay from "./TrackPlay.svelte";

    /** @type {import('$lib/api/api').CollectionTrack} */
    export let track;

    /** @type {number} */
    export let index;

</script>

<tr on:click={() => {}} draggable="true" on:drop={() => {}} on:dragover={() => {}} on:dragstart={() => {}}>
    <td><TrackPlay index={index} track_id={track.id}/></td>
    <td class="font-semibold">{track.name}</td>
    <td class="text-xs">{track.artist}</td>
    <td class="text-xs">{track.album_name}</td>
    <td class="text-xs">{convertDuration(track.duration)}</td>
    <td><FavouriteButton id={track.id} favourited={false}/></td>
</tr>

<style>
    td {
        @apply cursor-pointer;
        @apply select-none;
        text-overflow: ellipsis;
        white-space: nowrap;
    }

    tr {
        @apply transition-all;
        @apply text-gray-500;
    }

    tr:hover {
        @apply text-yellow-600;
    }
</style>