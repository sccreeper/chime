<script>
    import { onMount } from "svelte";
    import Main from "./Main.svelte";
    import Password from "./main/general/Password.svelte";
    import { session_object, user_object, view } from "./stores";
    import Cookies from "js-cookie"
    
    let username = "";
    let password = "";

    let login_status = "";
    let status_colour = "text-white";

    function login() {
        
        var data = new FormData()

        data.append("u", username)
        data.append("p", password)

        fetch("/api/auth", {
            method: "POST",
            body: data
        }).then((response) => response.json() ).then((data) => {

            if (data.status == "incorrect") {
                login_status = "Username or password incorrect."
                status_colour = "text-red-600"
            } else {
                login_status = "Login details correct. Loading..."
                status_colour = "text-lime-400"

                let date = new Date()
                date.setTime(date.getTime() + (7*24*60*60*1000))
                Cookies.set("session", btoa(JSON.stringify(data.session)).replaceAll("/", "-").replaceAll("+", "_").replaceAll("=", "."), {expires: date, path: "/", same_site: false})
                Cookies.set("user", btoa(JSON.stringify(data.user)).replaceAll("/", "-").replaceAll("+", "_").replaceAll("=", "."), {expires: date, path: "/", same_site: false})

                user_object.set(data.user)
                session_object.set(data.session)
                view.set(Main)
            
            }   

        })

    }

    onMount(() => {

        if (Cookies.get("session") != undefined && Cookies.get("user") != undefined) {
            
            let session_json = atob(Cookies.get("session").replaceAll("-", "/").replaceAll("_", "+").replaceAll(".", "="))
            let session_obj = JSON.parse(session_json)
            let user_json = atob(Cookies.get("user").replaceAll("-", "/").replaceAll("_", "+").replaceAll(".", "="))
            let user_obj = JSON.parse(user_json)
            user_object.set(user_obj)
            session_object.set(session_obj)
            view.set(Main)

        }

    })

</script>

<div id="login-div">
    <h1 class="text-lg p-1">Chime</h1>
    <br>
    <p class="text-xs p-1 {status_colour}">{login_status}</p>
    <br>
    <input type="text" bind:value={username} placeholder="Username"/>
    <br>
    <Password bind:value={password} placeholder="Password"/>
    <br>
    <button on:click={login} class="mt-3">Login</button>
</div>


<style>

    #login-div {
        position: absolute;
        left: 50%;
        top: 50%;
        transform: translate(-50%, -50%);
    }

</style>