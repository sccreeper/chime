<script>
    import MinorButton from "./general/MinorButton.svelte";
    import { openModal } from 'svelte-modals'
    import UploadModal from "./modals/UploadModal.svelte";
    import { active_view, search_results } from "../stores";
    import Settings from "./modals/Settings.svelte";

    let searchValue = ""
    let search_box;

    function executeSearch() {
        
        fetch("/api/search", {
            method: "POST",
            body: JSON.stringify({query: searchValue})
        }).then(resp => resp.json()).then(data => {
            search_results.set(data)

            active_view.set({name: "search", id: ""})
        })

    }

    function searchEnter(e) {
        if (e.code == "Enter") {
            executeSearch()
        }
    }

    function upload_modal() {
        openModal(UploadModal)
    }

    function settings() {
        openModal(Settings)
    }

    function profile() {
        
    }

</script>

<svelte:window />

<div class="grid-top flex-none m-2">

    <div class="w-3/4 grid grid-cols-2 grid-rows-1 grid-search">

        <input type="text" bind:value={searchValue} placeholder="Search" class="search" on:keypress={searchEnter}/>
        <MinorButton icon="search" callback={executeSearch} bind:this={search_box}/>

    </div>
    <div class="grid grid-cols-3 grid-rows-1 gap-4 items-center justify-items-center text-xl">
        <MinorButton icon="plus-lg" callback={upload_modal}/>
        <MinorButton icon="gear-fill" callback={settings}/>
        <MinorButton icon="person-circle" callback={profile}/>
    </div>
    

</div>

<style>

.grid-top {

    display: grid;
    columns: 2;
    grid-template-columns: 8fr 2fr;
    align-items: center;
    justify-items: center;

}

.grid-search {
    align-items: center;
    justify-items: center;
    grid-template-columns: 9fr 1fr;
}

.search {
    @apply w-full border-b-2 p-3 mr-4;
}
</style>
