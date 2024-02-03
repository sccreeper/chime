export interface Collection {
    title: string;
    cover: string;
    is_album: boolean;
    description: string;
    protected: boolean;

    tracks: CollectionTrack[]

    dates: string[]
}

export interface CollectionTrack {
    id: string;
    name: string;
    released: number;
    artist: string;
    album_id: string;
    duration: number;
    cover_id: string;
}