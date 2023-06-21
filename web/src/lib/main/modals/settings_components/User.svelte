<script>
    import { get_current_component } from "svelte/internal";
    import { edit_user_error, user_object } from "../../../stores";
    import { allowed_username_chars, verifyString } from "../../../util";
    import MinorButton from "../../general/MinorButton.svelte";
    import { get } from "svelte/store";

    export let username = "";
    export let is_admin = false;
    export let id = "";

    const comp = get_current_component();

    let password = "••••••••••";
    let processing_request = false;

    function editUsername() {
        let new_username = prompt("Enter a new username", username);

        if (new_username != null) {
            if (new_username == "") {
                return;
            } else if (!verifyString(username, allowed_username_chars)) {
                edit_user_error.set("Only allowed characters A-Z a-z 0-9 _-");
            } else {
                fetch("/api/admin/change_username", {
                    method: "POST",
                    body: JSON.stringify({
                        username: new_username,
                        password: "",
                        user_id: id,
                    }),
                }).then((resp) => {
                    if (resp.ok) {
                        username = new_username;
                    } else {
                        edit_user_error.set(
                            "There was an error changing the username."
                        );
                    }
                });
            }
        }
    }

    function deleteUser() {
        if (id == get(user_object).user_id) {
            edit_user_error.set("You cannot delete yourself.");
            return
        }

        let conf = confirm(`Are you sure you want to delete the user ${username}?\n
        This will permanently delete all of their assets (tracks, radios, playlists, albums, covers)\n
        This action *cannot be undone*.`);

        if (conf) {
            fetch("/api/admin/delete_user", {
                method: "POST",
                body: JSON.stringify({ user_id: id }),
            }).then((resp) => {
                if (resp.ok) {
                    comp.$destroy();
                } else {
                    edit_user_error.set(
                        "There was an error deleting the user."
                    );
                }
            });
        }
    }

    function resetPassword() {
        fetch("/api/admin/reset_password", {
            method: "POST",
            body: JSON.stringify({ user_id: id }),
        })
            .then((resp) => resp.json())
            .then((data) => {
                password = data.password;
            });
    }

    function toggleAdmin() {
        if (processing_request) { return }

        processing_request = true;

        fetch("/api/admin/toggle_admin", {
            method: "POST",
            body: JSON.stringify({
                user_id: id
            })
        }).then(
            resp => {
                if (resp.ok) {
                    processing_request = false;
                    is_admin = !is_admin;
                }
        })

    }

    
</script>

<tr>
    <td class="font-semibold">{username}<MinorButton hint="Edit username" callback={editUsername} icon="pencil"/></td>
    
    <td>{password}<MinorButton hint="Reset user password" callback={resetPassword} icon="arrow-clockwise"/></td>
    
    <td><i on:click={toggleAdmin} class="bi bi-{is_admin ? "check-lg" : "x-lg"} {processing_request ? "cursor-wait" : "cursor-pointer"} hover:text-yellow-500 transition-all"/></td>
    
    <td><MinorButton hint="Delete user" callback={deleteUser} icon="trash"/></td>
</tr>
