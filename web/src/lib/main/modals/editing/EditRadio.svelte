<script>
    import { closeModal } from "svelte-modals";
    import MinorButton from "../../general/MinorButton.svelte";
    import MinorButtonText from "../../general/MinorButtonText.svelte";
    import { active_view } from "../../../stores";
    import { createNotification, notificationID } from "../../../notifications";
    import Notification from "../../notifications/Notification.svelte";

    export let radio_title = ""
    export let radio_url = ""
    export let data_callback = function (data) {}
    
    let processing_request = false
    let error_text = ""

    function apply() {
      processing_request = true

      fetch("/api/edit/radio", {
        method: "POST",
        body: JSON.stringify({

            radio_id: $active_view.id,
            name: radio_title,
            url: radio_url

        })
      }).then(resp => {
        if (!resp.ok) {
            error_text = "There was an error."
        } else {
            data_callback({
                title: radio_title,
                url: radio_url
            })

            closeModal()
            createNotification(Notification, {id: notificationID(), text: "Changed radio details", icon:"check-lg", expiry: 5000})
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

        <input type="text" placeholder="Title" bind:value={radio_title}/>
        <input type="text" placeholder="URL" bind:value={radio_url}/>

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
  