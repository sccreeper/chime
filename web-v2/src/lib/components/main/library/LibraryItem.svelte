<script>

    import {page} from "$app/stores";

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

   $: href = `/app/main/${type}/${item_id}`

</script>

<a {href} class:current="{$page.url.pathname === href}">{name}</a>

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

    a.current {
        @apply text-yellow-600;
    }

    a:hover {
        @apply text-yellow-600;
    }

</style>