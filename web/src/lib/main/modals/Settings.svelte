<script>
    import { closeModal } from "svelte-modals";
    import MinorButton from "../general/MinorButton.svelte";
    import { user_object } from "../../stores";
    import { get } from "svelte/store";
    import Password from "../general/Password.svelte";
    import MinorButtonText from "../general/MinorButtonText.svelte";

    export let isOpen

    let old_password = ""
    let new_password0 = ""
    let new_password1 = ""
    let password_error = {text: "", ok: false}

    let username_password = ""
    let username_change = get(user_object).username
    let username_error = {text: "", ok: false}

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
            username_error = {text: "Username can only contain characters A-Z a-z _-", ok: false}
            break;
          case "bad_auth":
            username_error = {text: "Incorrect password", ok: false}
            break;
          default:
            username_error = {text: `Username changed to ${username_change}`, ok: true}
            break;
        }
      })

    }

    function changePassword() {
      if (old_password == "" || new_password0 == "" || new_password1 == "") {
        password_error = {text: "Passwords cannot be empty", ok: false}
      } else if (new_password0 != new_password1) {
        password_error = {text: "Passwords must match", ok: false}
      } else {

        fetch("/api/admin/change_password", {
          method: "POST",
          body: JSON.stringify({
            old_password: old_password,
            new_password_0: new_password0,
            new_password_1: new_password1,
          })
        }).then(response => {
          if (!response.ok) {
            password_error = {text: "There was an error changing the password", ok: false}
          } else {
            password_error = {text: "Password changed successfully", ok: true}
          }
        })

      }

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

      <p class="text-sm small-heading">Change username</p>
      <p class="text-xs {username_error.ok ? `text-green-400` : `text-red-400`}">{username_error.text}</p>
      <input type="text" bind:value={username_change} placeholder="Username"/>
      <Password bind:value={username_password} placeholder="Password"/>

      <MinorButtonText callback={changeUsername} icon="box-arrow-in-right" text="Change username"/>

      <p class="text-sm small-heading">Change password</p>
      <p class="text-xs {password_error.ok ? `text-green-400` : `text-red-400`}">{password_error.text}</p>
      <Password bind:value={old_password} placeholder="Old password"/>
      <Password bind:value={new_password0} placeholder="New password"/>
      <Password bind:value={new_password1} placeholder="Repeat new password"/>

      <MinorButtonText callback={changePassword} icon="box-arrow-in-right" text="Change password"/>
  
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