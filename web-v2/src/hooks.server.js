/** @type {import('@sveltejs/kit').HandleFetch} */
export async function handleFetch({request, fetch}) {
    
    if (!request.url.startsWith("http://chime:8042")) {
        
        request = new Request(
            "http://chime:8042" + request.url,
            request
        )

    }

    return fetch(request)

}