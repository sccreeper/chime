<script>
    import { closeModal } from "svelte-modals";
    import MinorButton from "../general/MinorButton.svelte";
    import { user_object } from "../../stores";
    import { get } from "svelte/store";

    export let isOpen

    let old_password = ""
    let new_password0 = ""
    let new_password1 = ""

    let username_password = ""
    let username_change = get(user_object).username
    let username_status = ""
    let username_good = false

    function changeUsername() {
      
      fetch("/api/admin/change_username", {
        method: "POST",
        body: JSON.stringify({
          username: username_change,
          password: username_password,
          user_id: get(user_object).user_id
        })
      }).then(response => response.json()).then(data => {
        switch (data.status) {
          case "bad_username":
            username_good = false
            username_status = "Username can only contain characters A-Z a-z _-"
            break;
          case "bad_auth":
            username_good = false
            username_status = "Incorrect password"
            break;
          default:
            username_good = true
            username_status = `Username changed to ${username_change}`
            break;
        }
      })

    }

    function changePassword() {
      
    }

</script>

{#if isOpen}
<div role="dialog" class="modal">
  <div class="contents">

      <span class="ml-auto"><MinorButton icon="x-lg" callback={closeModal}/></span>

      {#if $user_object.is_admin}
      <h1>Users</h1>
      {/if}
      <h1>Settings</h1>

      <p class="text-sm p-1">Change username</p>
      <p class="text-xs p-1 {username_good ? `text-green-400` : `text-red-400`}">{username_status}</p>
      <p class="text-xs p-1">Username</p>
      <input type="text" bind:value={username_change}/>
      <p class="text-xs p-1">Password</p>
      <input type="text" bind:value={username_password}/>

      <button on:click={changeUsername}>Change username</button>

      <p class="text-sm p-1">Change password</p>

      <p class="text-xs p-1">Old password</p>
      <input type="text" bind:value={old_password}/>
      <p class="text-xs p-1">New password</p>
      <input type="text" bind:value={new_password0}/>
      <p class="text-xs p-1">Repeat new password</p>
      <input type="text" bind:value={new_password1}/>

      <button on:click={changePassword}>Change Password</button>
  
  </div>
</div>
{/if}

<style>

.modal {
      position: fixed;
      top: 0;
      bottom: 0;
      right: 0;
      left: 0;
      display: flex;
      justify-content: center;
      align-items: center;
    }
  
    .contents {
      @apply bg-gray-700;
      
      display: flex;
      flex-direction: column;
      align-items: start;
      justify-content: center;

      min-width: 35vw;
      min-height: 30vh;
      border-radius: 6px;
      padding: 16px;
      display: flex;

    }
</style>