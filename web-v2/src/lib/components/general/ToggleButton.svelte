<script>
    export var icon = ""
    export var active = false

    /**
     * @callback toggleCallback
     * @param {boolean} active
    */

    /** @type {toggleCallback} */
    export var callback;

    function click() {
        active = !active;
        callback(active)
    }

</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<i on:click={click} class="bi bi-{icon} {active ? "active" : "inactive"}"></i>

<style lang="postcss">
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