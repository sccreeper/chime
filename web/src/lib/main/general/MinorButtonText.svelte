<script>
    export var icon = ""
    export var callback = function (e) {}
    export var text = ""

</script>

<span on:click={callback}><i class="bi bi-{icon}"></i> <span class="span-text">{text}</span></span>

<style>
span {
    @apply text-lg;
    @apply text-slate-400;
    @apply transition-all;
    cursor: pointer;
}

.span-text {
    @apply text-sm;
}

span:hover {
    @apply text-yellow-500;
}
</style>