<script>
    import HorizontalDivider from "../general/HorizontalDivider.svelte";
    import ListDivider from "../general/ListDivider.svelte";
    import MinorButtonText from "../general/MinorButtonText.svelte";

    /** @type {boolean} */
    export let showModal;
    export let title;
    export let icon;

    /**
     * @type {HTMLDialogElement}
     */
    let dialog;

    $: if (dialog && showModal) dialog.showModal()
    $: if (!showModal && dialog) dialog.close()

</script>

<dialog bind:this={dialog}>

    <ListDivider icon={icon} text={title} />

    <HorizontalDivider/>

    <slot name="body"/>

    <HorizontalDivider/>

    <div class="button-div">

    <MinorButtonText icon="x-lg" text="Close" callback={() => {dialog.close(); showModal = false;}}/>

    <slot name="buttons"/>

    </div>

</dialog>


<style lang="postcss">

    dialog {
        @apply bg-gray-700;
        @apply text-white;
        @apply p-3;
        @apply w-1/4;
    }

    dialog::backdrop {
        @apply bg-black;
        @apply bg-opacity-25;
    }

    dialog[open] {
        animation: fade 0.25s cubic-bezier(0.215, 0.610, 0.355, 1);
        @apply outline-none;
    }

    dialog[open]::backdrop {
        animation: fade 0.25s ease-out;
    }

    @keyframes fade {
        from {
            opacity: 0;
        }
        to {
            opacity: 1;
        }
    }

    .button-div {
        @apply flex flex-row gap-3;
        @apply m-2;
    }

</style>