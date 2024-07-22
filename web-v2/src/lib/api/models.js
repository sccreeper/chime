/**
 * @typedef {Object} Collection
 * @property {string} title
 * @property {string} cover
 * @property {boolean} is_album
 * @property {string} description
 * @property {boolean} protected
 * 
 * @property {CollectionTrack[]} tracks
 * @property {string[]} dates
 */

/**
 * @typedef {Object} CollectionTrack
 * @property {string} id
 * @property {name} string
 * @property {number} released
 * @property {string} artist
 * @property {string} album_id
 * @property {number} duration
 * @property {string} cover_id
 * @property {number} disc
 * @property {string} album_name
 */

/**
 * Contains the TrackMetadata for an object. Can be see as an extension of {CollectionTrack}
 * @typedef {Object} TrackMetadata
 * @property {string} title
 * @property {string} album_name
 * @property {string} cover_id
 * @property {string} artist
 * @property {string} original_file
 * @property {string} format
 * @property {number} duration
 * @property {number} released
 * @property {number} size
 * @property {number} position
 * @property {number} disc
 */

/** 
 * @typedef {object} SearchResults
 * @property {SearchTrack[]} tracks
 * @property {SearchCollection[]} collections
 * @property {SearchRadio[]} radios
 */

/**
 * @typedef {Object} SearchTrack
 * @property {string} id
 * @property {string} album_id
 * @property {string} artist
 * @property {string} title
 * @property {number} duration
 * @property {string} cover
 */

/**
 * @typedef {Object} SearchCollection
 * @property {string} id
 * @property {string} title
 * @property {string} cover
 * @property {boolean} is_album
 */

/**
 * @typedef {Object} SearchRadio
 * @property {string} id
 * @property {string} name
 * @property {string} cover
 */


/** 
 * @typedef {Object} Session
 * @property {string} sessionId
 * @property {string} userId
 * @property {boolean} isAdmin
 * @property {string} username
 */

/**
 * @typedef {Object} StorageInformation
 * @property {number} totalVolumeSpace
 * @property {number} usedByOthers
 * @property {number} usedByChime
 * @property {Object} breakdown
 * @property {number} breakdown.backups
 * @property {number} breakdown.cache
 * @property {number} breakdown.covers
 * @property {number} breakdown.tracks
 */

/**
 * @typedef {Object} User
 * @property {string} id
 * @property {string} username
 * @property {boolean} isAdmin
 */
export {}