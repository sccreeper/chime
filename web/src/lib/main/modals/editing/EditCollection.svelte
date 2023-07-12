<script>
  import { closeModal } from "svelte-modals";
  import MinorButton from "../../general/MinorButton.svelte";
  import MinorButtonText from "../../general/MinorButtonText.svelte";
  import { album_list } from "../../../stores";
  import { createNotification, notificationID } from "../../../notifications";
  import Notification from "../../notifications/Notification.svelte";

  function apply() {
    processing_request = true;

    fetch("/api/edit/collection", {
      method: "POST",
      body: JSON.stringify({
        collection_id: collection_id,
        name: collection_title,
        description: collection_description,
        is_album: collection_is_album,
      }),
    }).then((resp) => {
      // Handle if action was not performed
      if (!resp.ok) {
        error_text = "There was an error. Try again later.";
        processing_request = false;
      } else {
        processing_request = false;
        
        // Return inputted data to component that opened modal.
        data_callback({
          collection_description: collection_description,
          collection_title: collection_title,
        });
        
        // Update the album list with new names
        fetch("/api/get_collections", {
          method: "GET",
        }).then((response) => response.json())
        .then((data) => {
          album_list.set(data);
        });

        // Finally close the modal
        closeModal();
        createNotification(Notification, {id: notificationID(), text: "Changed collection details", icon:"check-lg", expiry: 5000})
      }
    });
  }

  export let collection_title = "";
  export let collection_description = "";
  export let collection_is_album = false;
  export let collection_id = "";
  export let data_callback = function callback(params) {};

  let error_text = "";
  let processing_request = false;

  // provided by Modals
  export let isOpen;
</script>

{#if isOpen}
  <div role="dialog" class="modal">
    <div class="contents">
      <span class="ml-auto"><MinorButton icon="x-lg" callback={closeModal} /></span>

      <p class="text-xs text-red-600">{error_text}</p>

      <input type="text" placeholder="Collection title" bind:value={collection_title}/>

      <input type="text" placeholder="Collection description" bind:value={collection_description}/>
      
      <label for="edit_is_album_check">Is album?</label>
      <input type="checkbox" id="edit_is_album_check" bind:checked={collection_is_album}/>

      <MinorButtonText text="Apply" icon="check-lg" callback={apply} bind:disabled={processing_request}/>
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
