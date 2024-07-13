<script>
    /** @type {string} */
    export let icon;
    /** @type {function(): void} */
    export let callback;
    /** @type {string} */
    export let text;
    /** @type {boolean} */
    export let disabled = false;
    /** @type {string} */
    let cursor_pointer;
    /** @type {string} */
    export let hint = "";

    $: cursor_pointer = disabled ? "not-allowed" : "pointer"

    function handleClick() {
        if (!disabled) {
            callback()
        }
    }

</script>

<p style="cursor: {cursor_pointer}" on:click={handleClick} title={hint}><i class="bi bi-{icon}"></i> <span class="span-text">{text}</span></p>

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