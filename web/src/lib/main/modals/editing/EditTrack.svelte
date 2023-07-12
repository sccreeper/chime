<script>
    import { closeModal } from "svelte-modals";
    import MinorButton from "../../general/MinorButton.svelte";
    import MinorButtonText from "../../general/MinorButtonText.svelte";
    import { album_list, track_metadata_view } from "../../../stores";
    import { get } from "svelte/store";
    import { createNotification, notificationID } from "../../../notifications";
    import Notification from "../../notifications/Notification.svelte";

    export let track_title = ""
    export let track_artist = ""
    export let track_album_id = ""
    export let track_released = ""
    export let data_callback = function (data) {}

    let previous_album_id = track_album_id

    let processing_request = false
    let error_text = ""

    function apply() {
      processing_request = true

      fetch("/api/edit/track", {
        method: "POST",
        body: JSON.stringify({
          track_id: get(track_metadata_view),
          name: track_title,
          released: parseInt(track_released, 10),
          artist: track_artist,
          album_id: track_album_id
        })
      }).then(resp => {
        if (resp.ok) {
          
          let track_album_name = ""

          // Perform callback and close modal

          processing_request = false
          data_callback({
            title: track_title,
            artist: track_artist,
            released: track_released,
            album_id: track_album_id,
            album_name: track_album_name,
            previous_album_id: previous_album_id,
          })

          closeModal()
          createNotification(Notification, {id: notificationID(), text: "Changed track details", icon:"check-lg", expiry: 5000})

        } else {
          error_text = "There was an error."
        }
      })
    }

    // provided by Modals
    export let isOpen;
  </script>
  
  {#if isOpen}
    <div role="dialog" class="modal">
      <div class="contents">
        <span class="ml-auto"><MinorButton icon="x-lg" callback={closeModal} /></span>
        
        <p class="text-xs text-red-600">{error_text}</p>

        <input type="text" placeholder="Title" bind:value={track_title}/>
        <input type="text" placeholder="Artist" bind:value={track_artist}/>
        <input type="text" placeholder="Released" accept="[0-9]" maxlength="4" bind:value={track_released}/>
        <select bind:value={track_album_id}>
            {#each $album_list.albums as a}
                <option value={a.id}>{a.name}</option>
            {/each}
        </select>

        <MinorButtonText text="Apply" icon="check-lg" bind:disabled={processing_request} callback={apply}/>

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
      @apply z-10;
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
  