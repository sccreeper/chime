<script>
    import { api_favourite_track } from "../../edit";
    import { favourites_list } from "../../stores";


    export let id = ""
    export let favourited = false;

    function toggle() {
        
        favourited != favourited
        api_favourite_track(id)

    }

</script>

<i on:click={toggle} class="{($favourites_list.includes(id) || favourited) ? 'bi bi-star-fill text-yellow-500' : 'bi bi-star text-gray-500'}" ></i>

<style>
    i {
        @apply transition-colors;
    }
</style>