<script>
    import { convertDuration } from "../../util";


    export let id = ""
    export let title = ""
    export let artist = ""
    export let duration = 0
    export let album_name = ""

</script>

<tr><td><strong>{title}</strong></td><td>{artist}</td><td>{album_name}</td><td>{convertDuration(duration)}</td></tr>

<style>

    td {
        white-space: nowrap;
    }

</style>