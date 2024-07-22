<script>
    import ListDivider from "$lib/components/general/ListDivider.svelte";

    /** @type {import("./$types").PageData} */
    export let data;

</script>
<ListDivider text="Users" icon="people-fill"/>
<p>Add, delete and modify users on this server.</p>

<h1>User list</h1>
<p>{data.users.length} user(s)</p>

{#each data.users as user}
    <p>{user.username}</p>
{/each}