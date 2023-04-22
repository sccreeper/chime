<script>
    import Main from "./Main.svelte";
    import { session_object, view } from "./stores";

    
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
                document.cookie =  `session=${JSON.stringify(data)}; expires=${date.toUTCString()}; path=/; SameSite=none; Secure`

                session_object.set(data)
                view.set(Main)
            
            }   

        })

    }

</script>

<div id="login-div">
    <h1 class="text-lg p-1">Chime</h1>
    <br>
    <p class="text-xs p-1 {status_colour}">{login_status}</p>
    <br>
    <p class="text-xs p-1">Username</p>
    <input type="text" bind:value={username}/>
    <br>
    <p class="text-xs p-1">Password</p>
    <input type="password" bind:value={password}/>
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