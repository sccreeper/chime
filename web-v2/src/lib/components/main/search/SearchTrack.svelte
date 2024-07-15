<script>
    import { goto } from "$app/navigation";


    import no_cover from "$lib/assets/no_cover.png"
    import { convertDuration } from "$lib/util";

    /** @type {import("$lib/api/models").SearchTrack} */
    export let data;

    $: realCoverSrc = data.cover == "0" ? no_cover : `/api/collection/get_cover/${data.cover}?width=32&height=32`

    function navigate() {
        goto(`/app/main/collection/${data.album_id}`)
    }

</script>

<tr on:click={navigate}>

    <td><img src={realCoverSrc} alt="Cover for {data.title}" width="16" height="16"></td>
    <td>{data.title}</td>
    <td class="text-xs whitespace-nowrap text-ellipsis block">{data.artist}</td>
    <td class="text-xs">{convertDuration(data.duration)}</td>

</tr>

<style lang="postcss">

    tr {
        @apply transition-all;
        @apply text-gray-500;
        @apply select-none;
        @apply cursor-pointer;
    }

    tr:hover {
        @apply text-yellow-600;
    }

    td > img {
        @apply border;
        @apply transition-all;
        @apply border-gray-500;
    }

    tr:hover td > img {
        @apply border-yellow-600;
    }

</style>