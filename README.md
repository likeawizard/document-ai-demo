# Google DocumentAI API Test

## Installation
* `cp config.dev.yml config.yml`
* Set `config.yml` values:
    * `app` `debug: true` will log some extra info
    * `secret` placeholder for secretive things the app might do - signing JWTs, etc...
    * `store` currently only supports `driver: os|gcloud` -
        * `os` stores files on the filesystem in the `location` folder. Make sure `location` exists
        * `gcloud` stores files in the `location` bucket on GCloud storage.
            * **TODO** Credentials files are hardcoded. Add options to specify a creds file specifically for storage or use a global GCLoud creds file between the processor and store
            * **TODO** The storage bucket is set to public. This is not a production ready solution and a major privacy breach. The bucket should be set to private and URLs should be created via the `SignedURL` method to not expose the data.
    * `database` `driver: inmemory|sqlite`
        * `inmemory` is a simple non-persistant store. Only usable for quick debugging and in tests to naively mock a database.
        * `sqlite` with `name: <filename.db>` is persistant but still only usable for dirty prototyping
        * `sqlite` can sometimes have `CGO` related compile issues. Make sure `CGO` is enabled or comment out the driver import in `database/sqlite.go`
    * `document-ai` `project-id`, `processor-id` and `location` should be set according to your processor endpoint. See: [DocumentAI Request](https://cloud.google.com/document-ai/docs/send-request#curl)
        * `credsfile` service account credentials. Make sure the service account has permissions for Document AI. See: [Google Cloud Service Accounts](https://developers.google.com/workspace/guides/create-credentials#service-account)

* Compile `go build -o api main.go` and run `./api` or simply run the app `go run .`

## Test it
* Run unit tests with `go test ./...`
* The app implements a REST API for `expenses/`
* Currently only `POST` and `GET` requests are supported.
* POST `expenses/`
    * Payload `Content-Type: multipart/form-data` with a single `file` field
    * Sample request with `curl`
        ```
        curl -X POST http://localhost:8080/expenses \
            -F "file=@receipt3.png" \
            -H "Content-Type: multipart/form-data"
        ```
    * On successful request returns a `json` response:
        ```
        {
            "id": "83bfe566-4254-4333-8ed1-7a54f918e796",
            "filename": "document5.pdf",
            "status": "pending",
            "mime_type": "application/pdf",
            "path": "83bfe566-4254-4333-8ed1-7a54f918e796.pdf",
            "json_path": "83bfe566-4254-4333-8ed1-7a54f918e796.json"
        }
        ```
* GET `expenses/{uuid}`
    * Sample request with `curl`
        ```
        curl http://localhost:8080/expenses/61b36905-5745-4167-8b6c-5e796445216a
        
        ```
    * On successful request returns a `json` response:
        ```
        {
            "id": "83bfe566-4254-4333-8ed1-7a54f918e796",
            "filename": "document5.pdf",
            "status": "ready",
            "mime_type": "application/pdf",
            "path": "83bfe566-4254-4333-8ed1-7a54f918e796.pdf",
            "json_path": "83bfe566-4254-4333-8ed1-7a54f918e796.json"
        }
        ```
        * `id` generated `UUID`
        * `filename` original upload filename
        * `status` the record status `pending`, `ready`, `failed`
            * When a new record is uploaded an entry with `pending` status is created
            * The document processor sends a request to the processor API and updates with a `ready` status on completion
            * Or if the request timeouts or any other error is encountered the status is set to `failed`
        * `mime_type` uploaded file MIME Type. See supported formats: [DocumentAI file types](https://cloud.google.com/document-ai/docs/file-types)
        * `path` and `json_path` are the stored filenames. Does not include the `store.location` directory. But using the same file store it will retrieve the file correctly. This should be reworked in a complete app and is only a demo version.
* **TODO** PATCH `expenses/` - allow adding `Tags` to processed expenses
    
## Improvements & Scalability
* While still only a simple Demo/Test app, I for the most part tried to make it as functional and clean as possible.
* Thread safety. The Gin framework and the use of `go routines` in the `GoogleDocumentAI.Process` method are almost guaranteed to cause panics related to concurrent reads/writes. For example the `inmemory` database implemented on `map[uiid]Record` will cause panics in when multiple request will be handled at the same time. It could be solved by using `mutex` or thread safe maps or other thread safe solutions.
* All services and clients are handled via an interface. Thus it should be relatively easy to swap out the FileSystem based file storage with a cloud based storage bucket or any other solution by simply creating an adapter with the appropriate wrapper that implements the interface.
* Document processing can take some time up to several seconds and maybe longer for larger documents. Currently the app creates a record of a request and initializes it with a `pending` status. While the user could refresh the API `expense/` endpoint with their record `UUID` it would be better that the initial request also adds a `callBackURI` that the client provides to recive a notification on completion or failure.
* In a real world application I believe it would be best to separate the upload and processing parts in separate micro services which would communicate via a message broker, pub/sub or any other method. The processor service could have a worker pool architecture and and subscribe to processing request messages. One could then spin up as many processors on say K8s to scale according to demand.
* It would also be a good practice to have separate data models for the `Record` type. One for internal database and storage and one for exposure to API endpoints. Keeping them separate adds more verbosity and some duplication of code but in a complex application having different types for the same data based on context often makes more sense than one-fits-all solution.
* I have no delusions of grandeur that my code will have no bugs, never crash or not have performance issues. So a production ready solution should have more robust logging, performance metric collection by `Grafana` and some issue/fault tracker with `Sentry` or similar solutions. To be able to monitor the performance and help fix issues.
    