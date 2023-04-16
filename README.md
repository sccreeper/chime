# chime

This repository is the monorepo for chime, a locally hosted music streaming server.

---

### Running/Building

It is recommended you run chime in Docker using docker-compose.
You can do this by running `./docker-run.sh`
This builds the project in Docker. It also creates a volume in your home directory for Chime's data.

---

### Project structure

- `server` - Contains the source code for the backend web server.
- `web` - Contains the source code for the web UI. This is served by the server.
- `app` - Source code for the desktop and mobile apps.