# FileService

## Setup

To start using this service you will need following .env variables

- `FILES_PATH` - Directory where files would be saved (defaults to srv/files)
- `FILE_SERVICE_PORT` - Port on which server will run (defaults to `8000`)
- `MAX_FILE_SIZE` - Maximum file size server will accept (defaults to 2Mib)

## Api

- `GET /api/files/{filename}` - Gets specified file
- `POST /api/files` - Posts multipart file (key should be `file`)