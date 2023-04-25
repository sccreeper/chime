<script>
    import { closeModal } from 'svelte-modals'
    import HorizontalDivider from '../general/HorizontalDivider.svelte'
    import MinorButton from '../general/MinorButton.svelte';
    import { get } from 'svelte/store';
    import { album_list } from '../../stores';

    function extractCollections() {
        
        let arr = []

        get(album_list).albums.forEach(element => {
            arr.push(element)
        });

        get(album_list).playlists.forEach(element => {
            arr.push(element)
        });

        return arr

    }

    function add(id) {
        
        fetch("/api/collection/add_track", {
            method: "POST",
            body: JSON.stringify({
                track_id: track_id,
                collection_id: id,
            })
        }).then(() => {
            closeModal()
        })

    }
    
    export let track_id = ""

    // provided by Modals
    export let isOpen
  
  </script>
  
  {#if isOpen}
  <div role="dialog" class="modal">
    <div class="contents">

        <span class="ml-auto"><MinorButton icon="x-lg" callback={closeModal}/></span>

        <p>Add to playlist</p>

        <HorizontalDivider/>

        <div class="h-32 w-80 overflow-y-scroll">

            {#each extractCollections() as item}
                
                <p on:click={() => {add(item.id)}} class="text-sm text-gray-500 cursor-pointer hover:text-yellow-600 transition-all">{item.name}</p>

            {/each}

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