<script>
    /** @type {import('./$types').LayoutData} */
    export let data;
</script>

<svelte:head>
    <title>Home - Chime</title>
</svelte:head>

<h1>Main</h1>

{#each data.lib.albums as item}
    <p>{item.name}</p>
{/each}

<slot/>