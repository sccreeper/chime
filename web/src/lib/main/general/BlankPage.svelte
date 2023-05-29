<script>

    export var icon = ""
    export var text = ""

</script>

<div class="container">

<div class="child">

<i class="bi bi-{icon}"/>
<h1>{text}</h1>

</div>

</div>

<style>

.container {
    @apply w-full;
    @apply h-full;

    @apply flex;
    @apply items-center;
    @apply justify-center;
}

.child {

    @apply flex;
    @apply flex-col;
    @apply text-center;

}

h1 {

    @apply font-semibold;
    @apply text-2xl;
    @apply text-gray-600;
    @apply select-none;

}

i {
    @apply text-9xl;
    @apply text-gray-600;
    @apply mb-4;
}

</style>