<script>
    import MinorButton from "./general/MinorButton.svelte";

    let searchValue = ""

    function executeSearch() {
        
    }

    function upload() {
        
    }

    function settings() {
        
    }

    function profile() {
        
    }

</script>

<div class="grid-top flex-none">

    <div class="w-3/4 grid grid-cols-2 grid-rows-1 grid-search">

        <input type="text" bind:value={searchValue} placeholder="Search" class="appearance-none w-full bg-transparent border-b-gray-500 border-b-2 text-gray-400 placeholder:text-gray-400 outline-none mr-4"/>
        <MinorButton icon="search" callback={executeSearch}/>

    </div>
    <div class="grid grid-cols-3 grid-rows-1 gap-2 items-center justify-items-center text-xl">
        <MinorButton icon="upload" callback={upload}/>
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

</style>
