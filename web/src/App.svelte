<script>
    import Login from "./lib/Login.svelte";
    import { view } from "./lib/stores";

  console.log("Hello World from Chime!")

  view.set(Login)
</script>

<main>
  <svelte:component this={$view}></svelte:component>
</main>


