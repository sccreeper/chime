<script>
    import { closeModal } from 'svelte-modals'
    import HorizontalDivider from '../general/HorizontalDivider.svelte'
    import MinorButton from '../general/MinorButton.svelte';


    export let message = ""
    export let callback = function (confirmed) {}

    // provided by Modals
    export let isOpen
  
  </script>
  
  {#if isOpen}
  <div role="dialog" class="modal">
    <div class="contents">

        <span class="ml-auto"><MinorButton icon="x-lg" callback={closeModal}/></span>

        <p>{message}</p>

        <HorizontalDivider/>

        <div class="w-48 flex flex-row items-center gap-3">

            <button on:click={() => {callback(true)}}>Yes</button>
            <button on:click={() => {callback(false)}}>No</button>

        </div>

    </div>
  </div>
  {/if}
  
  <style>
    .modal {
      position: fixed;
      top: 0;
      bottom: 0;
      right: 0;
      left: 0;
      display: flex;
      justify-content: center;
      align-items: center;
    }
  
    .contents {
      @apply bg-gray-700;
      
      display: flex;
      flex-direction: column;
      align-items: start;
      justify-content: center;

      min-width: 240px;
      border-radius: 6px;
      padding: 16px;
      display: flex;

    }
  
  </style>