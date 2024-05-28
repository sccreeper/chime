<script>
    import { convertDuration } from "$lib/util";
    import { getContext } from "svelte";
    import FavouriteButton from "./FavouriteButton.svelte";
    import TrackPlay from "./TrackPlay.svelte";
    import { page } from "$app/stores";
    import { PLAYER_CONTEXT_KEY } from "$lib/player";

    /** @type {import('$lib/player').ChimePlayer} */
    const {collectionId, playing, currentTrack} = getContext(PLAYER_CONTEXT_KEY)

    /** @type {import('$lib/api/api').CollectionTrack} */
    export let track;

    /** @type {number} */
    export let index;

    /** @type {number} */
    export let displayed_index;

</script>

<tr on:click={() => {}} draggable="true" on:drop={() => {}} on:dragover={() => {}} on:dragstart={() => {}} class:playing={$page.data.collection_id == $collectionId && $playing && $currentTrack?.id == track.id}>
    <td><TrackPlay index={index} track_id={track.id} displayed_index={displayed_index}/></td>
    <td class="font-semibold">{track.name}</td>
    <td class="text-xs">{track.artist}</td>
    <td class="text-xs">{track.album_name}</td>
    <td class="text-xs">{convertDuration(track.duration)}</td>
    <td><FavouriteButton id={track.id} favourited={false}/></td>
</tr>

<style lang="postcss">
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

    tr:hover, tr.playing {
        @apply text-yellow-600;
    }
</style>