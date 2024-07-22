<script>
    import ListDivider from "$lib/components/general/ListDivider.svelte";
    import { enhance } from "$app/forms";

    /** @type {import("./$types").ActionData} */
    export let form;

</script>
<ListDivider text="Security" icon="shield-fill"/>
<p>Change password & username</p>

<h1>Change password</h1>

<form action="?/changePassword" method="post" use:enhance>

<label>
    Old password <br>
    <input type="password" name="old"/>
</label>
<br>
<label>
    New password <br>
    <input type="password" name="new"/>
</label>
<br>
<label>
    Confirm new password <br>
    <input type="password" name="newConfirm"/>
</label>

<br>

{#if form?.success && form}
    <p class="text-green-400 font-semibold">Password changed successfully</p>
{:else if !form?.success && form}
    <p class="text-red-400 font-semibold">{form?.message}</p>
{/if}

<button type="submit" class="mt-3">Change password</button>

</form>

<h1>Change username</h1>

<form action="?/changeUsername" method="post"></form>