<script>
    /** @type {string} */
    export let icon;

    export let callback = () => {};
    /** @type {string} */
    export let hint = ""

</script>

<i on:click={callback} class="bi bi-{icon}" title={hint}></i>

<style>
i {
    @apply text-lg;
    @apply text-slate-400;
    @apply transition-all;
    cursor: pointer;
}

i:hover {
    @apply text-yellow-500;
}
</style>