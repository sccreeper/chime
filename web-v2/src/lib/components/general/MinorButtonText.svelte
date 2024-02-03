<script>
    /** @type {string} */
    export var icon;
    /** @type {function(): void} */
    export var callback;
    /** @type {string} */
    export var text;
    /** @type {boolean} */
    export var disabled = false;
    /** @type {string} */
    let cursor_pointer;

    $: cursor_pointer = disabled ? "not-allowed" : "pointer"

    function handleClick() {
        if (!disabled) {
            callback()
        }
    }

</script>

<p style="cursor: {cursor_pointer}" on:click={handleClick}><i class="bi bi-{icon}"></i> <span class="span-text">{text}</span></p>

<style lang="postcss">
p {
    @apply text-lg;
    @apply text-slate-400;
    @apply transition-all;

}

p:hover {
    @apply text-yellow-500;
}

.span-text {
    @apply text-sm;
}

</style>