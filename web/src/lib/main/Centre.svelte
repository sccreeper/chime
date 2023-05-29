<script>
    import { active_view } from "../stores";
    import Search from "./Search.svelte";
    import AlbumList from "./centre_components/AlbumList.svelte";
    import AlbumView from "./centre_components/AlbumView.svelte";
    import DetailsView from "./centre_components/DetailsView.svelte";
    import RadioView from "./centre_components/RadioView.svelte";
    import BlankPage from "./general/BlankPage.svelte";

</script>
<div class="flex h-full grow overflow-hidden">
    <div class="left">
        <AlbumList/>
    </div>

    <div class="grow centre">

        {#if $active_view.name === ""}
        <div class="flex justify-items-center items-center w-full h-full">

            <BlankPage icon="eye" text="Not viewing anything."/>
        
        </div>
        {:else}

        {#if $active_view.name === "radio"}
        <RadioView/>
        {:else if $active_view.name === "album"}
        <AlbumView/>
        {:else if $active_view.name === "search"}
        <Search/>
        {/if}

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
