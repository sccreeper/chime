# API documentation <!-- omit in toc -->
---
Contains all information about all API endpoints.

---

- [General info](#general-info)
- [Misc](#misc)
  - [/api/ping](#apiping)
- [Uploading/Downloading/Adding](#uploadingdownloadingadding)
  - [/api/upload](#apiupload)
  - [/api/add\_radio](#apiadd_radio)
  - [/api/download\_original/:track\_id](#apidownload_originaltrack_id)
- [Collections](#collections)
  - [/api/get\_collection/:collection\_id](#apiget_collectioncollection_id)
  - [/api/get\_collections](#apiget_collections)
- [Tracks](#tracks)
  - [/api/get\_track\_metadata/:track\_id](#apiget_track_metadatatrack_id)


## General info

- A response code of `200` indicates that the request was completed successfully.
- `403` The user completing this request is either not logged in or does not have permissions to modify/access the requested resource.
- `400` Something is wrong with the data you sent to the server.
- `500` Something has gone wrong on the server side and is generally not your fault.

## Misc

### /api/ping
- Method: **GET**
- Info: Is server alive?
- Returns: 
```json
{"message":"pong"}
```

## Uploading/Downloading/Adding

### /api/upload
- Method: **POST**
- Info: Upload a single track.
- Accepts: Multipart Form with field `file` for uploaded file. File can be MP3, FLAC, OGG or WAV.
- Returns: `200`

### /api/add_radio
- Method: **POST**
- Info: Add a radio
- Accepts:
```json
{
    "name:":"Example", //Any string
    "url":"https://example.com/radio/stream.m3u8" //Or any other internet radio stream
}
```
- Returns: `200`
### /api/download_original/:track_id
- Method: **GET**
- Info: Download the original file for a track.
- Accepts: 
  - `:track_id` - URL paramater, hexadecimal representation of Track ID (`int64`).
- Returns: The track file.

**For adding collections (playlists/albums) see ()[]**

## Collections

### /api/get_collection/:collection_id
- Method: **GET**
- Info: Gets metdata and tracks contained within collection.
- Accepts: Hexadecimal ID of collection in `collection_id` paramater.
- Response:

```json
{
  "title":"The Dark Side Of The Moon",
  "cover":"677cbefe80f855a0", //If the Cover ID is 0, there is no cover for this collection and a placeholder cover should be displayed instead.
  "is_album":true, //False for playlists
  "description":"One of the best albums from Pink Floyd",
  "protected":false, //If true, it can't be modified by any user.
  "tracks": [
    {
      "id":"7e4e87d3749a1c7b", //Hexadecimal ID of track
      "name":"Speak To Me",
      "released":1973,
      "artist":"Pink Floyd",
      "album_id":"ebd2c57e0e94860",
      "duration":64,
      "cover_id":"677cbefe80f855a0",
    }
    //...
  ],
  "dates":["1973"], //For an album, this is the date that it was released, for a playlist it is a list of dates when tracks were added to the playlist.

}
```

### /api/get_collections

**Note: In the future this will be renamed to `/api/user/get_library`.**

- Method: **GET**
- Info: Gets a the Names and IDs for all high level items in a user's library. This includes: Albums, Playlists and Radios.
- Accepts: No data is sent to the server.
- Response:

```json
{
  "albums":[
    {
      "id":"3548b6bc21c4fb89",
      "name":"The Age of Consent"
    }
    // ...
  ],
  "playlists":[
    {
      "id":"677cbefe80f855a0",
      "name":"All Songs"
    }
    // ...
  ],
  "radios":[
    {
      "id":"3bcb9baca172caf2",
      "name":"Radio 2"
    }
    // ...
  ]
}
```

## Tracks

### /api/get_track_metadata/:track_id

- Method: **GET**
- Info: Get the metadata about a particular track.
- Accepts: Hexadecimal ID of track in `track_id` parameter.
- Response:

```json
{
  "title":"Your Song",
  "album_name":"Elton John",
  "album_id":"66476bcbed4c002d",
  "cover_id":"5fdc00a07b746064",
  "artist":"Elton John",
  "original_file":"01-01-Paul_Buckmaster-Your_Song-SMR.flac", //The name of the file when this track was originally uploaded.
  "format":"FLAC", //The audio format for this track.
  "duration":244, //Duration in seconds
  "released":2013,
  "size":88138694 //The size of the track in bytes.
}
```
