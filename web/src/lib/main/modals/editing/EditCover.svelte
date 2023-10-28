<script>
    import { closeModal } from "svelte-modals";
    import MinorButton from "../../general/MinorButton.svelte";
    import { onMount } from "svelte";
    import CoverCard from "./CoverCard.svelte";
    import MinorButtonText from "../../general/MinorButtonText.svelte";
    import { get } from "svelte/store";
    import { active_view } from "../../../stores";
  
    export let current_cover = "";
    export let target_cover = "";
    export let data_callback = (id) => {};

    let cover_ids = [];

    onMount(() => {
      
      fetch("/api/library/get_covers").then(resp => resp.json()).then(data => {
      
        // Filter cover IDs

        let cover_ids_temp = []

        data.cover_ids.forEach(element => {
          
          if (element != current_cover) {
            cover_ids_temp.push(element)
          }

        });

        cover_ids = cover_ids_temp;

      })

    })

    function uploadCover() {

      console.log("lol")

      var cover_form = document.createElement("input")
      cover_form.setAttribute("type", "file");

      cover_form.addEventListener("change", async () => {

        if (cover_form.files.length > 0) {
          
          let file = cover_form.files[0]

          let data = new FormData();
          data.append("file", file)

          let resp = await fetch("/api/upload_cover", {method: "POST", body: data})
          let resp_data = await resp.json()

          cover_ids = [...cover_ids, resp_data.id]

        }

      })

      cover_form.click();
      
    }

    // provided by Modals
    export let isOpen;
  </script>
  
  {#if isOpen}
    <div role="dialog" class="modal">
      <div class="contents">
        <span class="ml-auto"><MinorButton icon="x-lg" callback={closeModal} /></span>
  
        <h1><i class="bi bi-image"></i> Pick cover</h1>

        <MinorButtonText icon="upload" callback={uploadCover} text="Upload new cover"/>

        <div class="cover-grid">
          
          {#if current_cover == "0"}
            <CoverCard id="0" selected="true" target_id={get(active_view).id} data_callback={data_callback}/>
          {:else}
            <CoverCard id="0" selected="false"/>
            <CoverCard id="{current_cover}" target_id={get(active_view).id} selected="true" data_callback={data_callback}/>
          {/if}

          {#each cover_ids as cover}
            <CoverCard id="{cover}" target_id={get(active_view).id} selected="false" data_callback={data_callback}/>
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

    .cover-grid {

      @apply grid;
      @apply grid-cols-5;

    }

    h1 {
      @apply text-yellow-600;
    }

  </style>