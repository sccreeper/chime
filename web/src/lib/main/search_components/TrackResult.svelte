<script>
    import { active_view } from "../../stores";
    import { convertDuration } from "../../util";


    export let title = ""
    export let cover = ""
    export let album_id = ""
    export let duration = 0.0
    export let artist = ""

    function setView() {
        active_view.set({
            name: "album",
            id: album_id,
        })
    }

</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<tr on:click={setView} class="display_div">
<td style="text-align: center;"><img src="/api/collection/get_cover/{cover}?width=32&height=32" alt="Album cover for song {title}" width="32" height="32"/></td>
<td class="font-semibold text-base">{title}</td>
<td>{artist}</td>
<td>{convertDuration(duration)}</td>
</tr>

<style>
    td {
        @apply cursor-pointer;
        @apply select-none;
        text-overflow: ellipsis;
        white-space: nowrap;
        @apply text-gray-500;
        @apply text-xs;
    }
</style>