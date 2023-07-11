<script>
    import { closeModal } from 'svelte-modals'
    import HorizontalDivider from '../general/HorizontalDivider.svelte'
    import MinorButton from '../general/MinorButton.svelte';
    import MinorButtonText from '../general/MinorButtonText.svelte';
    import { createNotification, notificationID, removeNotification } from '../../notifications';
    import UploadNotification from '../notifications/UploadNotification.svelte';
    import Notification from '../notifications/Notification.svelte';
    
    let radio_name = ""
    let radio_url = ""
    let error_message = ""

    let collection_name = ""
    let collection_description = ""
    let collection_cover = null
    let collection_is_album = false

    // Upload multiple or single files

    function uploadFiles() {

        let file_form = document.createElement("input")
        file_form.setAttribute("type", "file")
        file_form.setAttribute("multiple", "")

        file_form.addEventListener("change", async () => {

            let files = file_form.files;
            let notification_id;
            let failed = false;
            
            for (let i = 0; i < files.length; i++) {
                notification_id = notificationID()
                createNotification(UploadNotification, {id: notification_id, progress: i+1, count: files.length, finished: false})

                const element = files[i];
                
                let data = new FormData();
                data.append("file", element)

                var resp = await fetch("/api/upload", {method: "POST", body: data})
                if (resp.ok) {
                  removeNotification(notification_id)
                  console.log("Uploaded track successfully!")
                } else {
                  removeNotification(notification_id)
                  createNotification(Notification, {text: `Error with: ${element.name}`, icon: "x-lg", expires: 5000})
                  break
                }

            }

            if (!failed) {
              removeNotification(notification_id)
              notification_id = notificationID()
              createNotification(UploadNotification, {id: notification_id, progress: files.length, count: files.length, finished: true})
            }

        })

        file_form.click()

    }
    
    // Add radio method

    function addRadio() {
        
        fetch("/api/add_radio", {
            method: "POST",
            body: JSON.stringify({
                name: radio_name,
                url: radio_url
            })
        }).then(() => {

            closeModal()

        })

    }

    // Add collection

    function addCollection() {
      
      let data = new FormData()

      let use_custom_cover

      if (collection_cover == null) {
        use_custom_cover = false
      } else {
        use_custom_cover = true
      }

      data.append("data", JSON.stringify({

        name: collection_name,
        description: collection_description,
        is_album: collection_is_album,
        custom_cover: use_custom_cover,

      }))

      if (collection_cover != null) {
        data.append("cover", collection_cover) 
      }

      fetch("/api/collection/add", {
            method: "POST",
            body: data
        }).then(() => {

            closeModal()

        })

    }

    function setCoverFile() {
      
      let input = document.createElement("input")
      input.setAttribute("type", "file")
      input.setAttribute("accept", "image/png, image/jpeg, image/webp")

      input.addEventListener("change", () => {

        collection_cover = input.files[0]

      })

      input.click()

    }

    // provided by Modals
    export let isOpen
  
  </script>
  
  {#if isOpen}
  <div role="dialog" class="modal">
    <div class="contents">

        <span class="ml-auto"><MinorButton icon="x-lg" callback={closeModal}/></span>

        <h1>Upload file</h1>
        <MinorButtonText icon="upload" text="Add files" callback={uploadFiles}/>

        <span class="w-full"><HorizontalDivider/></span>

        <h1>Add radio</h1>

        <input type="text" bind:value={radio_name} placeholder="Name"/>
        <input type="text" bind:value={radio_url} placeholder="URL"/>
        <p class="text-red-600 text-xs p-1">{error_message}</p>
        <MinorButtonText text="Add" icon="plus-lg" callback={addRadio}/>

        <span class="w-full"><HorizontalDivider/></span>

        <h1>Add collection</h1>

        <input type="text" bind:value={collection_name} placeholder="Name"/>
        <input type="text" bind:value={collection_description} placeholder="Description"/>
        <p class="text-xs p-1">Cover (optional)</p>
        <MinorButtonText icon="upload" text="Add image" callback={setCoverFile}/>
        <label for="is_album_check" class="text-xs">Is album?</label>
        <input type="checkbox" id="is_album_check" bind:checked={collection_is_album}/>

        <MinorButtonText icon="plus-lg" text="Add" callback={addCollection}/>

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
  
    .actions {
      margin-top: 32px;
      display: flex;
      justify-content: flex-end;
    }
  </style>