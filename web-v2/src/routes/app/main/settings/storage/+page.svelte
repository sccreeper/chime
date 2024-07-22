<script>
    import StackedBar from "$lib/components/general/graphs/StackedBar.svelte";
    import ListDivider from "$lib/components/general/ListDivider.svelte";

    /** @type {import("./$types").PageData} */
    export let data;

    /**
     * Exists to avoid repearting that massive one liner inside it.
     * @param {number} val
     * @param {number} pow
     * @returns {number}
     */
    const humaniseStorageVal = (val, pow) => {
        if (val == 0) {
            return val
        } else {
            return parseFloat((val / Math.pow(10, pow)).toPrecision(4))
        }
    }

    /** @type {import("$lib/components/general/graphs/types").ChartData[]}*/
    let chartData = [
        {
            label: "Used by Chime",
            value: humaniseStorageVal(data.storageData.usedByChime, 9),
            colour: "#facc15",
            units: "GB"
        },
        {
            label: "Used by others",
            value: humaniseStorageVal(data.storageData.usedByOthers, 9),
            colour: "#52525b",
            units: "GB"
        },
        {
            label: "Unused",
            value: humaniseStorageVal(data.storageData.totalVolumeSpace-(data.storageData.usedByChime+data.storageData.usedByOthers), 9),
            colour: "#27272a",
            units: "GB"
        }
    ]

    /** @type {import("$lib/components/general/graphs/types").ChartData[]}*/
    let breakdownChartData = [
        {
            label: "Tracks",
            value: humaniseStorageVal(data.storageData.breakdown.tracks, 6),
            colour: "#ca8a04",
            units: "MB"
        },

        {
            label: "Covers",
            value: humaniseStorageVal(data.storageData.breakdown.covers, 6),
            colour: "#facc15",
            units: "MB"
        },

        {
            label: "Cache",
            value: humaniseStorageVal(data.storageData.breakdown.cache, 6),
            colour: "#fef08a",
            units: "MB"
        },

        {
            label: "Backups",
            value: humaniseStorageVal(data.storageData.breakdown.backups, 6),
            colour: "#713f12",
            units: "MB"
        },
    ]

</script>

<ListDivider text="Storage" icon="device-hdd-fill"/>
<p>Monitor the amount of storage being used by Chime.</p>

<h1>Outline</h1>

<StackedBar data={chartData}/>

<p><b>Total:</b> {parseFloat((data.storageData.totalVolumeSpace / Math.pow(10, 9)).toPrecision(4))}GB</p>

<h1>Breakdown</h1>

<StackedBar data={breakdownChartData}/>