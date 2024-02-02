<script>
    /**
     * @type {string}
     */
    export let item_id;
    /**
     * @type {string}
     */
    export let name;
    /**
     * @type {string}
    */
   export let type;

</script>

<a href="/app/main/{type}/{item_id}">{name}</a>

<style lang="postcss">

    a {
        @apply text-gray-300;
        @apply transition-all;
        cursor: pointer;
        @apply overflow-ellipsis;
        @apply overflow-hidden;
        @apply whitespace-nowrap;
        @apply block;
    }

    a:hover {
        @apply text-yellow-600;
    }

    .active {
        @apply text-yellow-600;
    }

</style>