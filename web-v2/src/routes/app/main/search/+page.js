
/** @type {import('./$types').PageLoad} */
export async function load({url, fetch}) {

    const params = url.searchParams
    /** @type {string} */
    const query = params.get("query") ?? ""

    if (query == "") {

        return {searchResults: {tracks: [], collections: [], radios: []}}
        
    } else {
        const req = await fetch(
            `/api/search`,
            {
                method: "POST",
                body: JSON.stringify({query: query})
            }
        )

        /** @type {import("$lib/api/models").SearchResults} */
        const res = await req.json()

        return {searchResults: res, query: query}
    }

}