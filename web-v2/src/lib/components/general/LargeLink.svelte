<script>

    import { page } from "$app/stores";

    /** @type {string} */
    export let text;
    /** @type {string} */
    export let icon;

    /** @type {string} */
    export let href;

</script>

<a {href} class:current={$page.url.pathname === href}>

    <i class="bi bi-{icon}"></i>
    {text}

</a>

<style lang="postcss">

    a {

        @apply text-xl;
        @apply text-gray-300;
        @apply transition-all;
        display: inherit;
        
        &:hover, &.current {
            @apply text-yellow-600;
        }

    }

</style>