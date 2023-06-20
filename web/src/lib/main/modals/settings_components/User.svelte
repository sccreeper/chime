<script>
    import MinorButton from "../../general/MinorButton.svelte";

    export let username = ""
    export let is_admin = false
    export let id = ""

    function editUser() {
        
    }

    function deleteUser() {
        
    }


</script>

<tr>
    <td class="font-semibold">{username}</td>
    <td>••••••••••</td>
    <td>{#if is_admin}<i class="bi bi-check-lg"></i>{:else}<i class="bi bi-x-lg"></i>{/if}</td>
    <td><MinorButton callback={editUser} icon="pencil"></MinorButton></td>
    <td><td><MinorButton callback={deleteUser} icon="trash"></MinorButton></td>
</tr>