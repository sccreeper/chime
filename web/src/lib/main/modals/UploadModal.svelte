<script>
    import { closeModal } from 'svelte-modals'
    import HorizontalDivider from '../general/HorizontalDivider.svelte'
    import MinorButton from '../general/MinorButton.svelte';
    import MinorButtonText from '../general/MinorButtonText.svelte';
    
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
            
            for (let i = 0; i < files.length; i++) {
                const element = files[i];
                
                let data = new FormData();
                data.append("file", element)

                var resp = await fetch("/api/upload", {method: "POST", body: data})
                if (resp.ok) {
                    console.log("Uploaded track successfully!")
                } else {
                    console.log("There was an error uploading the track!")
                }

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

        <p>Upload file</p>
        <button on:click={uploadFiles} class="mt-2">Add files</button>

        <span class="w-full"><HorizontalDivider/></span>

        <p>Add radio</p>

        <p class="text-xs p-1">Name</p>
        <input type="text" bind:value={radio_name}/>
        <p class="text-xs p-1">URL</p>
        <input type="text" bind:value={radio_url}/>
        <p class="text-red-600 text-xs p-1">{error_message}</p>
        <button on:click={addRadio} class="mt-2">Add</button>

        <span class="w-full"><HorizontalDivider/></span>

        <p>Add collection</p>

        <p class="text-xs p-1">Name</p>
        <input type="text" bind:value={collection_name}/>
        <p class="text-xs p-1">Description</p>
        <input type="text" bind:value={collection_description}/>
        <p class="text-xs p-1">Cover (optional)</p>
        <MinorButtonText icon="upload" text="Add image" callback={setCoverFile}/>
        <p class="text-xs p-1">Is album?</p>
        <input type="checkbox" bind:checked={collection_is_album}/>

        <button on:click={addCollection} class="mt-2">Add</button>

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
  
    .actions {
      margin-top: 32px;
      display: flex;
      justify-content: flex-end;
    }
  </style>