<script>
    export var icon = ""
    export var active = false
    export var callback = function (active) {}

    function click() {
        active = !active;
        callback(active)
    }

</script>

<i on:click={click} class="bi bi-{icon} {active ? "active" : "inactive"}"></i>

<style>
i {
    @apply text-lg;
    @apply cursor-pointer;
}

.inactive {
    @apply text-slate-400;
}

.active {
    @apply text-yellow-500;
}
</style>