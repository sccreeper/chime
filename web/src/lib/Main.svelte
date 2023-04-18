<script>
    import Centre from "./main/Centre.svelte";
    import Player from "./main/Player.svelte";
    import TopBar from "./main/TopBar.svelte";

    
</script>

<div class="flex flex-col h-full">
    <TopBar/>
    <Centre/>
    <Player/>
</div>

<style>

    .main {
        grid-template-rows: 1fr 8fr 1fr;
        width: 100%;
        height: 100%;
    }

</style>
