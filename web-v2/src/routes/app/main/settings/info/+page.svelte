<script>
import ListDivider from "$lib/components/general/ListDivider.svelte";

/** @type {import("./$types").PageData} */
export let data;

</script>
<ListDivider text="Server information" icon="info-circle-fill"/>
<p>Information about the server</p>

<p><b>Server ID:</b> {data.serverId}</p>