<script>
    export var icon = ""
    export var callback = function (e) {}

</script>

<i on:click={callback} class="bi bi-{icon}"></i>

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