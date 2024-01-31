<img src="./web/src/assets/logo.png" width="256" height="256"/>

---

This repository is the monorepo for chime, a locally hosted music streaming server.

---

### Note

This project is only a "small" hobby project of mine. For personal use I use [Jellyfin](https://github.com/jellyfin/jellyfin) as even though Chime has a large list of features and is comparable to the music aspects of Jellyfin, developing it and trying to make it work for multiple platforms properly is too big a load for a solo developer like myself especially when I want to work on other projects. As a result commits to this repo will probably be few and far between from now on as I focus on other small projects which are more manageable.

TLDR; Making a fully-featured music streaming server is too much, so this will just be added to as and when I get the chance.

### About

Hobby music streaming server made as an experiment.

- Server is written in Go
- Web interface written in (bad) Svelte<sup>1<sup/>.
- Mobile app with download support written in Dart & Flutter.
- Automatically parses music metadata when added to your library.
- Internet radio support (Uses [hls.js](https://github.com/video-dev/hls.js/)).
- Partial Chromecast support.

### Running/Building

It is recommended you run chime in Docker using docker-compose.
You can do this by running `./docker-run.sh`
This builds the project in Docker. It also creates a volume in your home directory for Chime's data.

---

### Project structure

- `server` - Contains the source code for the backend web server.
- `web` - Contains the source code for the web UI. This is served by the server.
- `app` - Source code for the desktop and mobile apps.
- `castproxy` - Source code for the Chromecast proxy to enable control of Chromecast devices from platforms where the Cast SDK is not supported.
- `web-v2` - A rewrite of the existing frontend using SvelteKit instead of standard Svelte.

### Credits

- Bootstrap Icons
- Svelte
- Gin
- Flutter
- just_audio
- get_it
- [Chromecast Cast Button](https://icons8.com/icon/1I0NE97niMwR/chromecast-cast-button) icon by [Icons8](https://icons8.com)