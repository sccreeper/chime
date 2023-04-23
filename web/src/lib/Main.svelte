<script>
    import Centre from "./main/Centre.svelte";
    import Player from "./main/Player.svelte";
    import TopBar from "./main/TopBar.svelte";

    import { Modals, closeModal } from 'svelte-modals'

    
</script>


<Modals>
    <div
      slot="backdrop"
      class="backdrop"
      on:click={closeModal}
    />
</Modals>
  

<div class="flex flex-col h-full">
    <TopBar/>
    <Centre/>
    <Player/>
</div>

<style>

    .backdrop {
    position: fixed;
    top: 0;
    bottom: 0;
    right: 0;
    left: 0;
    background: rgba(0,0,0,0.50)
  }

</style>
