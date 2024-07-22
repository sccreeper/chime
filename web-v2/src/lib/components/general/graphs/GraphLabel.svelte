<script>

/** @type {import("$lib/components/general/graphs/types").ChartData}*/
export let data;

</script>

<p><span style:background-color={data.colour}></span> <b>{data.label}:</b> {data.value}{data.units ?? ""}</p>

<style lang="postcss">

p > span {
    @apply w-3;
    @apply h-3;
    @apply inline-block;
    @apply border;
    @apply border-gray-600;
}

</style>