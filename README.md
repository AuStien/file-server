# File Server

A simple Go program to serve files, with some bonus features :sunglasses:

*NOTE: Only tested on Linux*

## Key points

- Serves all static files from a root folder
- Has a `GET /api/random` endpoint which returns the path to a file with whitelisted file extension
- Watches the filesystem to keep track of all currently existing files (no restart needed when new files/folders are added or old removed)

## Running

**Required** environment variables:

- `FILE_SERVER_ROOT_FOLDER`: A path to the folder containing files
- `FILE_SERVER_HOST`: The host where the server is running (will typically be `http://localhost`, or a domain: `https://example.com`)

*Optional* environment variables:

- `FILE_SERVER_WHITELISTED_EXTENSIONS`: A comma separated list to set custom file types to be considered valid. Defaults to `jpeg,jpg,png,gif,mp4,mp3,wav,avi,mkv,mov,webm`
- `PORT`: Set the port the server will be running on. Defaults to `8000`
- `LOG_LEVEL`: Set the desired log level, only `debug` and `info` are currently supported. Defaults to `info`

