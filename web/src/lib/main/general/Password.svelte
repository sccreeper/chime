<script>
    import Toggle from "../player_components/Toggle.svelte";

    export let value = ""
    export let placeholder = ""
    let shown = false

    function input(e) {
        value = e.target.value
    }

</script>


<div class="inline">
<input type={shown ? "text" : "password"} on:input={input} placeholder={placeholder}/>
<Toggle icon="eye" bind:active={shown}/>
</div>