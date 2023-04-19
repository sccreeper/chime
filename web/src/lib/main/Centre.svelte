<script>
    import AlbumList from "./centre_components/AlbumList.svelte";

</script>
<div class="flex h-full grow overflow-hidden">
    <div class="left">
        <AlbumList/>
    </div>

    <div class="grow centre">

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
