<script>
    import { onMount } from "svelte";
    import { Cast } from "../../cast";


    export let type = ""
    export let name = ""
    export let model = ""
    export let uuid = ""

    let icon = ""

    onMount(() => {
        switch (type) {
            case "audio":
                icon = "speaker-fill"
                break;

            case "cast":
                icon = "tv-fill"
                break;

            case "group":
                icon = "collection-play-fill"
                break;
            default:
                icon = "question-octagon-fill"
                break;
        }
    })

</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<div class="cast" on:click={() => Cast.connect(uuid)}>

    <i class="bi bi-{icon} text-2xl"></i>

    <div>
        <p class="name">{name}</p>
        <p class="model">{model}</p>
    </div>

</div>

<style>
    
    .cast {
        @apply flex flex-row items-center gap-2;
        @apply text-slate-400;
        @apply transition-all;
        @apply select-none cursor-pointer;
    }

    .cast:hover {
        @apply text-yellow-600;
    }

    .name {
        @apply text-sm;
    }

    .model {
        @apply text-xs;
    }

</style>