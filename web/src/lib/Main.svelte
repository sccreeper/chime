<script>
    import Centre from "./main/Centre.svelte";
    import Player from "./main/Player.svelte";
    import TopBar from "./main/TopBar.svelte";

    import { Modals, closeModal } from 'svelte-modals'
    import { createNotification, notifications } from "./notifications";
    import { onMount } from "svelte";
    import Notification from "./main/notifications/Notification.svelte";
    import { active_view } from "./stores";

    onMount(() => {
      createNotification(Notification, {icon: "check-lg", text: "Logged in", expiry: 3000})
    })
    
</script>


<Modals>
    <div
      slot="backdrop"
      class="backdrop"
      on:click={closeModal}
    />
</Modals>

<div class="fixed right-3 bottom-0 w-56 h-72 flex justify-end flex-col">
  {#each $notifications as n}
      <svelte:component this={n.component} {...n.props}/>    
  {/each}
</div>
  

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
    @apply backdrop-brightness-50;
    @apply transition;
  }

</style>
