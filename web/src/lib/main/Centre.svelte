<script>
    import { current_album } from "../stores";
import AlbumList from "./centre_components/AlbumList.svelte";
    import AlbumView from "./centre_components/AlbumView.svelte";

</script>
<div class="flex h-full grow overflow-hidden">
    <div class="left">
        <AlbumList/>
    </div>

    <div class="grow centre">

        {#if ($current_album == 0)}
        <div class="flex justify-items-center items-center w-full h-full">

            <h3>Not viewing anything. Click on an album or playlist to get started.</h3>
        
        </div>
        {:else}
        <AlbumView/>
        {/if}


    </div>

    <div class="right">

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

        min-width: 450px;
        background-color: rgb(45, 53, 66);

    }

</style>
