<script>
    export var icon = ""
    export var callback = function (e) {}
    export var text = ""

</script>

<p on:click={callback}><i class="bi bi-{icon}"></i> <span class="span-text">{text}</span></p>

<style>
p {
    @apply text-lg;
    @apply text-slate-400;
    @apply transition-all;
    cursor: pointer;
    

}

p:hover {
    @apply text-yellow-500;
}

.span-text {
    @apply text-sm;
}

</style>