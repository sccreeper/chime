<script>
	/** @type {import('./$types').PageData} */
	export let data;
</script>

<svelte:head>
    <title>Login - Chime</title>
</svelte:head>

<div>
    <h1 class="text-lg p-1">Chime</h1>
    <br>
    <form method="POST">
        <input name="username" type="text" placeholder="Username"/>
        <br>
        <input name="password" type="password" placeholder="Password"/>
        <br>
        <button class="mt-3">Login</button>
    </form>
</div>

<p class="text-center text-xs text-gray-400 bottom-0 absolute p-2">{data.server_id}</p>

<style>
    div {
        position: absolute;
        left: 50%;
        top: 50%;
        transform: translate(-50%, -50%);
    }
</style>

