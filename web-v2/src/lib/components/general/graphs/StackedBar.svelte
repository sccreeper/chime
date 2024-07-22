<script>
    import GraphLabel from "./GraphLabel.svelte";
    import { Canvas, Layer } from "svelte-canvas";

    /** @type {import("$lib/components/general/graphs/types").ChartData[]} */
    export let data;

    /** @type {string} */
    export let backgroundColour = "#000000"

    /** @type {import("svelte-canvas").Render} */
    const render = ({context, width, height}) => {

        context.fillStyle = backgroundColour;
        context.fillRect(0, 0, width, height)

        let total = data.reduce((accumulator, current) => {return accumulator + current.value}, 0)

        let totalWidth = 0
        data.forEach((value) => {

            context.fillStyle = value.colour
            context.fillRect(totalWidth, 0, (value.value/total)*width, height)

            totalWidth += (value.value/total)*width

        })

    }

</script>

{#each data as label}

<GraphLabel data={label}/>
    
{/each}

<div class="h-9 w-48">

<Canvas autoplay>

    <Layer {render}/>

</Canvas>

</div>
