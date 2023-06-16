<script>
    import Login from "./lib/Login.svelte";
    import { view } from "./lib/stores";
    import favicon from "./assets/logo.svg"

  console.log("Hello World from Chime!")

  view.set(Login)
</script>

<svelte:head>
  <link rel="icon" type="image/svg" href={favicon}/>
</svelte:head>

<svelte:component this={$view}></svelte:component>


