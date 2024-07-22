<script>
    import HorizontalDivider from "$lib/components/general/HorizontalDivider.svelte";
import ListDivider from "$lib/components/general/ListDivider.svelte";

    /** @type {import("./$types").LayoutData} */
    export let data;

</script>
<div class="flex flex-row h-full">

    <div class="w-1/5 flex-shrink border-r border-gray-600">

        <div class="m-2 sidebar">

            <ListDivider icon="gear-fill" text="Settings"/>

            <a href="/app/main/settings/info">Info</a>
            <a href="/app/main/settings/security">Security</a>

            {#if data.user.isAdmin}
                <ListDivider icon="shield-lock-fill" text="Admin"/>

                <a href="/app/main/settings/users">Users</a>
                <a href="/app/main/settings/storage">Storage</a>
                <a href="/app/main/settings/backups">Backups</a>
            {/if}

        </div>
    

    </div>

    <div class="flex-grow">

        <div class="m-2">

            <slot></slot>

        </div>

    </div>

</div>

<style lang="postcss">

    .sidebar a {
        @apply block;
    }

</style>