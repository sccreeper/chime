<script>
    import { active_view } from "../stores";
    import AlbumList from "./centre_components/AlbumList.svelte";
    import AlbumView from "./centre_components/AlbumView.svelte";
    import DetailsView from "./centre_components/DetailsView.svelte";
    import RadioView from "./centre_components/RadioView.svelte";

    let current_view = null;

    active_view.subscribe((val) => {

        if (val.name == "radio") {
            current_view = RadioView
        } else if (val.name == "album") {
            current_view = AlbumView
        }

    })

</script>
<div class="flex h-full grow overflow-hidden">
    <div class="left">
        <AlbumList/>
    </div>

    <div class="grow centre">

        {#if $active_view.name == ""}
        <div class="flex justify-items-center items-center w-full h-full">

            <h3>Not viewing anything. Click on an album or playlist to get started.</h3>
        
        </div>
        {:else}
        <svelte:component this={current_view}/>
        {/if}


    </div>

    <div class="right">
        <DetailsView/>
    </div>
</div>

<style>

    .centre {
        background-color: rgb(38, 45, 56);
    }

    .left {
        width: 15%;
        background-color: rgb(45, 53, 66);
    }

    .right {

        width: 450px;
        background-color: rgb(45, 53, 66);

    }

</style>
