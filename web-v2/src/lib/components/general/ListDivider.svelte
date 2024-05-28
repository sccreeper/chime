<script>
    /**
     * @type {string}
     */
    export let icon;
    /**
     * @type {string}
     */
    export let text;
</script>

<strong><i class="bi bi-{icon}"></i> {text} <slot></slot></strong>

<style lang="postcss">

strong {
    @apply text-yellow-500;
    @apply block;
} 

</style>