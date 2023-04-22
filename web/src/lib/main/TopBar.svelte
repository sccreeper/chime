<script>
    import { get } from "svelte/store";
    import { session_object } from "../stores";
    import MinorButton from "./general/MinorButton.svelte";

    let searchValue = ""

    function executeSearch() {
        
    }

    function upload() {
        
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

    function settings() {
        
    }

    function profile() {
        
    }

</script>

<div class="grid-top flex-none m-2">

    <div class="w-3/4 grid grid-cols-2 grid-rows-1 grid-search">

        <input type="text" bind:value={searchValue} placeholder="Search" class="appearance-none w-full bg-transparent border-b-gray-500 border-b-2 text-gray-400 placeholder:text-gray-400 outline-none mr-4"/>
        <MinorButton icon="search" callback={executeSearch}/>

    </div>
    <div class="grid grid-cols-3 grid-rows-1 gap-4 items-center justify-items-center text-xl">
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
