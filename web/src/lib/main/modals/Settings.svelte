<script>
    import { closeModal } from "svelte-modals";
    import MinorButton from "../general/MinorButton.svelte";
    import { edit_user_error, user_object } from "../../stores";
    import { get } from "svelte/store";
    import Password from "../general/Password.svelte";
    import MinorButtonText from "../general/MinorButtonText.svelte";
    import { onMount } from "svelte";
    import User from "./settings_components/User.svelte";
    import { allowed_username_chars, verifyString } from "../../util";

    export let isOpen

    let old_password = ""
    let new_password0 = ""
    let new_password1 = ""
    let password_error = {text: "", ok: false}

    let username_password = ""
    let username_change = get(user_object).username
    let username_error = {text: "", ok: false}

    let total_storage = 0
    let used_by_others = 0
    let used_by_chime = 0
    
    let show_new_user_ui = false
    let new_user_error = ""

    let new_user = {
      username: "",
      password: "",
      is_admin: false
    }

    let users = [{
      username: "",
      is_admin: false,
      id: "",
    }]

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
          case "username_exists":
            username_error = {text: "Username already exists", ok: false}
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

    async function addUser() {
      
      // Verify form data

      if (!verifyString(new_user.username, allowed_username_chars)) {
        new_user_error = "Only allowed characters A-Z a-z 0-9 _-"
        return
      }

      const new_user_request = await fetch("/api/admin/add_user", 
      {
        method: "POST", 
        body: JSON.stringify(new_user)
      })

      if (new_user_request.ok) {
        
        fetch("/api/admin/users", {method: "GET"}).then(resp => resp.json()).then(data => {
          users = data.users
        })

        show_new_user_ui = false;

      } else {

        new_user_error = "There was an error adding the new user"

      }

    }

    // When mounted fetch all required data, (this only applies if a user is admin)
    onMount(() => {

      if (get(user_object).is_admin) {      
        fetch("/api/admin/storage", 
        {
          method: "GET"
        }
        ).then(response => response.json()).then(
          data => {

            total_storage = data.total_volume_space
            used_by_others = data.used_by_others
            used_by_chime = data.used_by_chime
            
          }
        )
        
        fetch("/api/admin/users", {method: "GET"}).then(resp => resp.json()).then(data => {
          users = data.users
        })
        
      }

      edit_user_error.set("")

    })

</script>

{#if isOpen}
<div role="dialog" class="modal">
  <div class="contents">

      <span class="ml-auto"><MinorButton icon="x-lg" callback={closeModal}/></span>

      {#if $user_object.is_admin}
      <h1>Users</h1>
      {#if !show_new_user_ui}

      <MinorButtonText text="Add user" icon="plus-lg" callback={() => {show_new_user_ui = true}}/>
      
      {:else}

      <p class="text-xs text-red-600">{new_user_error}</p>
      <input type="text" placeholder="Username" bind:value={new_user.username}/>
      <input type="text" placeholder="Password" bind:value={new_user.password}/>
      <input type="checkbox" id="new_user_admin" bind:checked={new_user.is_admin}/>
      <label for="new_user_admin">Admin</label>
      <MinorButtonText text="Add" icon="plus-lg" callback={addUser}/>

      {/if}
      <p class="text-xs text-red-600">{$edit_user_error}</p>
      <table class="w-full text-gray-400">
        <colgroup>
          <col span="1" style="width: 45%;">
          <col span="1" style="width: 45%;">
          <col span="1" style="width: 10%;">
          <col span="1" style="width: 10%;">
        </colgroup>

        <tr class="font-light text-sm small-heading">
          <td>Username</td>
          <td>Password</td>
          <td>Admin</td>
          <td>Delete</td>
        </tr>
      {#each users as u}
        <User username={u.username} id={u.id} is_admin={u.is_admin}/>
      {/each}
      </table>  

      <h1>Storage</h1>
      <p>Total storage: {Math.round(total_storage / Math.pow(10, 9))}GB</p>
      <p>Used by other data: {Math.round(used_by_others / Math.pow(10, 9))}GB</p>
      <p>Used by Chime: {Math.round(used_by_chime / Math.pow(10, 6))}MB</p>

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
      @apply z-10;
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