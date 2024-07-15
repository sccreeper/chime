<script>

    /** @type {import("$lib/api/models").SearchCollection} */
    export let data;
    import noCover from "$lib/assets/no_cover.png";

    $: realCover = data.cover == "0" ? noCover : `/api/collection/get_cover/${data.cover}?width=160&height=160`

</script>

    <a href="/app/main/collection/{data.id}">
        <div>

            <img src={realCover} alt="Cover image for {data.title}"/>
            <p class="title">{data.title}</p>
            <p class="text-xs">{data.is_album ? "Album" : "Collection"}</p>

        </div>
    </a>

<style lang="postcss">

    div {
        @apply m-3;
        @apply p-3;
        @apply text-gray-300;
        @apply transition-all;
        @apply w-36;
    }

    div:hover {
        @apply bg-gray-700;
    }

    div > img {
        @apply border;
        @apply border-gray-500;
    }

    div:hover > img {
        @apply border-yellow-600;
    }

    div p:first-of-type {

        @apply whitespace-nowrap;
        @apply overflow-hidden;
        @apply overflow-ellipsis;
        @apply block;

    }


</style>