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

    // Password modificatin
    let old_password = ""
    let new_password0 = ""
    let new_password1 = ""
    let password_error = {text: "", ok: false}

    // Userame modification
    let username_password = ""
    let username_change = get(user_object).username
    let username_error = {text: "", ok: false}

    // Storage
    let storage_data = {
      total_volume_space: 0,
      used_by_others: 0,
      used_by_chime: 0,

      breakdown: {
        backups: 0,
        cache: 0,
        covers: 0,
        tracks: 0,
      }
    }

    // New users
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

    // Backup
    let backup_status = ""
    let backup_processing = false

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

    async function backup() {
      backup_processing = true

      let backup_id;
      let backup_interval;

      backup_status = "Starting backup..."
      
      let resp = await fetch("/api/admin/start_backup", {
        method: "GET",
      })
      let data = await resp.json()

      backup_id = data.id

      // Timeout to query status

      backup_interval = async () => {

        let resp = await fetch(`/api/admin/backup_status/${backup_id}`, {
          method: "GET"
        })
        let data = await resp.json()

        if (data.finished) {

          backup_processing = false

          backup_status = "Backup finished"

          let date = new Date()

          let link = document.createElement("a")
          link.download = `chime-${data.hash.substring(0, 16)}-${date.toDateString().replaceAll(" ", "")}.tar.gz`
          link.href = `/api/admin/download_backup/${backup_id}`
          link.click()
          document.removeChild(link)
          
        } else if (data.failed) {

          backup_status = "Backup failed"
          backup_processing = false

        } else {

          backup_status = `Backup progress: ${data.progress}%`
          setTimeout(backup_interval, 1000)

        }

      }

      setTimeout(backup_interval, 1000)

      
    }

    function clearBackups() {
      
      backup_status = "Removing backups..."
      backup_processing = true

      fetch("/api/admin/clear_backups", {
        method: "GET"
      }).then(resp => {
        if (resp.ok) {
          backup_status = "Backups removed sucessfully"
          backup_processing = false
        } else {
          backup_status = "There was an error removing the backups"
          backup_processing = false
        }
      })

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

            storage_data = data;
            
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
      <h1><i class="bi bi-people-fill"></i> Users</h1>
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

      <h1><i class="bi bi-device-hdd-fill"></i> Storage</h1>
      
      <p class="small-heading text-sm">System</p>
      <p>Total storage: {Math.round(storage_data.total_volume_space / Math.pow(10, 9))}GB</p>
      <p>Used by other data: {Math.round(storage_data.used_by_others / Math.pow(10, 9))}GB</p>
      <p>Used by Chime: {Math.round(storage_data.used_by_chime / Math.pow(10, 6))}MB</p>

      <p class="small-heading text-sm">Chime</p>
      <p>Backups: {Math.round(storage_data.breakdown.backups / Math.pow(10, 6))}MB</p>
      <p>Cache: {Math.round(storage_data.breakdown.cache / Math.pow(10, 3))}KB</p>
      
      <p>Tracks: {Math.round(storage_data.breakdown.tracks / Math.pow(10, 6))}MB</p>
      <p>Covers: {Math.round(storage_data.breakdown.covers / Math.pow(10, 6))}MB</p>

      <h1><i class="bi bi-database-fill-down"></i> Backup</h1>
      <MinorButtonText text="Download backup" icon="download" callback={backup} bind:disabled={backup_processing}/>
      <MinorButtonText text="Clear backups" icon="trash-fill" callback={clearBackups} bind:disabled={backup_processing}/>
      <p>{backup_status}</p>


      {/if}
      <h1><i class="bi bi-gear-fill"></i> Settings</h1>

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

h1 {
  @apply text-yellow-500;
}

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