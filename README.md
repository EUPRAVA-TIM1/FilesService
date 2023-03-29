# FileService

## Setup

To start using this service you will need following .env variables

- `FILES_PATH` - Directory where files would be saved (defaults to srv/files)
- `FILE_SERVICE_PORT` - Port on which server will run (defaults to `8000`)
- `MAX_IMAGE_SIZE` - Maximum image size server will accept (defaults to 2Mib)
- `MAX_PDF_SIZE` - Maximum pdf size server will accept (defaults to 30kb)

## Api

- `GET /api/files/{filename}` - Gets specified file
- `POST /api/files` - Posts multipart file (key should be `file`)