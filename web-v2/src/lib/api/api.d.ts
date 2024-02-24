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
    disc: number;
    album_name: string;
}

export interface TrackMetadata {
    title: string;
    album_name: string;
    album_id: string;
    cover_id: string;
    artist: string;
    original_file: string;
    format: string;
    duration: number;
    released: number;
    size: number;
    position: number;
    disc: number;
}