<script>
    import MinorButton from "./general/MinorButton.svelte";
    import { openModal } from 'svelte-modals'
    import UploadModal from "./modals/UploadModal.svelte";
    import { active_view, search_results } from "../stores";
    import Settings from "./modals/Settings.svelte";
    import Spinner from "./general/Spinner.svelte";
    import CastDevice from "./cast_components/CastDevice.svelte";
    import { Cast, cast_devices, current_cast_device, using_cast } from "../cast";
    import chromecast_icon from "../../assets/chromecast_icon.png";

    let searchValue = ""
    let search_box;

    let cast_hidden = true;

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

    function open_cast_menu() {

        cast_hidden = !cast_hidden;

        if (!cast_hidden) {
            console.log("Fetching cast devices...")

            cast_devices.set([])
            Cast.discover()
        }

    }

</script>

<div class="grid-top flex-none m-2">

    <div class="w-3/4 grid grid-cols-2 grid-rows-1 grid-search">

        <input type="text" bind:value={searchValue} placeholder="Search" class="search" on:keypress={searchEnter}/>
        <MinorButton icon="search" callback={executeSearch} bind:this={search_box}/>

    </div>
    <div class="grid grid-cols-4 grid-rows-1 gap-4 items-center justify-items-center text-xl">
        <MinorButton icon="plus-lg" callback={upload_modal} hint="Add"/>
        <MinorButton icon="gear-fill" callback={settings} hint="Settings"/>
        <!-- svelte-ignore a11y-click-events-have-key-events -->
        <div>
            <i on:click={open_cast_menu} class="bi bi-cast" title="Cast"></i>
            <div class="absolute w-48 h-40 overflow-y-scroll bg-gray-700 shadow-md -translate-x-24 p-2 {cast_hidden ? "hidden" : ""}">
                
                {#if $cast_devices.length == 0}
                <div class="w-full h-full flex items-center justify-center"><Spinner/></div>
                {:else}
                
                {#each $cast_devices as device}

                    <CastDevice name={device.name} model={device.model} type={device.type} uuid={device.uuid}/>
                    
                {/each}

                {/if}

            </div>
        </div>
    </div>

            
    {#if $using_cast}
    <p class="text-xs text-slate-400"><img class="inline" src={chromecast_icon} width="12" height="12" alt="chromecast icon"/> Connected to {$current_cast_device.name}</p>
    {/if}

</div>

<style>

.grid-top {

    display: grid;
    columns: 3;
    grid-template-columns: 8fr 2fr 1fr;
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

i {
    @apply text-lg;
    @apply text-slate-400;
    @apply transition-all;
    cursor: pointer;
}

i:hover {
    @apply text-yellow-500;
}
</style>
