<script>
    import { get } from "svelte/store";
    import { duration, playing, position, volume } from "../player";
    import Toggle from "./player_components/Toggle.svelte";
    import { convertDuration } from "../util";

    function skipForward(params) {
        
    }

    function skipBack(params) {
        
    }

    function toggleShuffle(state) {
        
    }

    function toggleRepeat(state) {
        
    }

    function mute(params) {
        
    }

    function volumeIcon(value) {

        console.log(value)

        if (value > 0.5) {
            return "volume-up-fill"
        } else if (value > 0.25) {
            return "volume-down-fill"
        } else if (value > 0) {
            return "volume-off-fill"
        } else {
            return "volume-mute-fill"
        }

    }

</script>

<div class="grid grid-rows-1 player-grid justify-items-center items-center flex-none m-2">

    <div class="grid grid-rows-2">

        <!-- Controls -->

        <div class="grid grid-rows-1 grid-cols-5 gap-3 items-center justify-items-center">
            <Toggle callback={toggleShuffle} icon="shuffle"/>
            <button class="control-button" on:click={skipBack}><i class="bi bi-skip-backward"></i></button>
            
            <button 
                class="control-button" 
                on:click={() => {playing.set(!get(playing))}}>
                    {#if !$playing}<i class="bi bi-play-fill"></i>{:else}<i class="bi bi-pause-fill"></i>{/if}
            </button>
            
            <button class="control-button" on:click={skipForward}><i class="bi bi-skip-forward"></i></button>
            <Toggle callback={toggleRepeat} icon="repeat"/>
        </div>

        <!-- Volume -->

        <div class="grid grid-rows-1 grid-cols-2 volume-grid items-center justify-items-center">
            <i class="bi bi-{volumeIcon($volume)} mr-1 text-slate-400" on:click={mute}></i><input class="seek" type="range" bind:value={$volume} min="0" max="1.0" step="0.05"/>
        </div>


    </div>

    <!-- Seek -->

    <div class="grid grid-cols-3 grid-rows-1 seek-grid w-full h-full items-center justify-items-center">
        <p class="duration select-none">{convertDuration($position)}</p>
        <input class="seek w-full" type="range" max={$duration} bind:value={$position}/>
        <p class="duration select-none">{convertDuration($duration)}</p>
    </div>

</div>

<style>

    .duration {
        @apply text-slate-400;
    }

    .player-grid {
        grid-template-columns: 2fr 8fr;
    }

    .volume-grid {
        grid-template-columns: 1fr 9fr;
    }

    .seek-grid {
        grid-template-columns: 1fr 4fr 1fr;
    }

    .seek {
        -webkit-appearance: none;
        background-color: transparent;
    }

    .seek::-webkit-slider-runnable-track {
        -webkit-appearance: none;
        appearance: none;

        background-color: transparent;
        border: none;
        box-shadow: none;

        @apply h-1;
        @apply bg-slate-500;
        @apply outline-none;
    }

    .seek::-moz-range-track {
        appearance: none;
        background-color: transparent;
        border: none;
        box-shadow: none;

        @apply h-1;
        @apply bg-slate-500;
        @apply outline-none;
    }


    .seek::-moz-range-thumb {
        -webkit-appearance: none;
        appearance: none;
        @apply h-3;
        @apply w-3;
        @apply bg-yellow-600;
        @apply outline-none;
        @apply border-none;
        cursor: pointer;
    }


    .seek::-webkit-slider-thumb {
        -webkit-appearance: none;
        @apply h-3;
        @apply w-3;
        @apply bg-yellow-600;
        @apply outline-none;
        @apply rounded-full;
        transform: translateY(-25%);
        cursor: pointer;
    }

    .control-button {

        border-radius: 50%;
        cursor: pointer;
        @apply w-10;
        @apply h-10;
        @apply text-lg;
    }



</style>