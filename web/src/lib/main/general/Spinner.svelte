
<!-- Spinning progress indicator -->

<div></div>

<style>

div {
    @apply border-4 border-gray-500;
    @apply border-t-4 border-t-yellow-600;
    @apply border-solid;
    border-radius: 50%;
    @apply w-7 h-7;
    animation: spin 1s linear infinite;
}

@keyframes spin {

    0% {transform: rotate(0deg);}
    100% {transform: rotate(360deg);}

}

</style>