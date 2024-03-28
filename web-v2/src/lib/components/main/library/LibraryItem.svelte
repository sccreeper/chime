<script>

    import {page} from "$app/stores";
    import noCover from "$lib/assets/no_cover.png";
    import { coverSizes } from "$lib/util";

    /**
     * @type {import('$lib/api/library').LibraryItem}
     */
    export let data;
    /**
     * @type {string}
    */
   export let type;

   $: href = `/app/main/${type}/${data.id}`
   $: src = data.cover_id == "0" ? noCover : `/api/collection/get_cover/${data.cover_id}?width=${coverSizes.icon}&height=${coverSizes.icon}`

</script>

<a {href} class:current="{$page.url.pathname === href}">
    <img {src} width="16" height="16"/> {data.name}
</a>

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

    a.current > img {
        @apply border-yellow-600;
        @apply border;
    }

    a:hover {
        @apply text-yellow-600;
    }

    a:hover > img {
        @apply border-yellow-600;
        @apply border;
    }

    a > img {
        @apply border-gray-500;
        @apply border;

        @apply inline;
    }


</style>