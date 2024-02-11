<script>
    export let id = ""
    export let favourited = false;

    function toggle() {

    }

</script>

<i on:click={toggle} class="{(false || favourited) ? 'bi bi-star-fill text-yellow-500' : 'bi bi-star text-gray-500'}" ></i>

<style>
    i {
        @apply transition-colors;
    }
</style>