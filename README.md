# Google DocumentAI API Test

## Installation
* `cp config.dev.yml config.yml`
* Set `config.yml` values:
    * `app` `debug: true` will log some extra info
    * `processor-driver` which service should be used for processing receipts
        * Could make a list of supported/active processors and more than one could be used to process a single receipt to improve data extraction via redundancy and second opinion.
    * `secret` placeholder for secretive things the app might do - signing JWTs, etc...
    * `store` currently only supports `driver: os|gcloud` -
        * `os` stores files on the filesystem in the `location` folder. Make sure `location` exists
        * `gcloud` stores files in the `location` bucket on GCloud storage.
            * **TODO** Credentials files are hardcoded. Add options to specify a creds file specifically for storage or use a global GCLoud creds file between the processor and store
            * **TODO** The storage bucket is set to public. This is not a production ready solution and a major privacy breach. The bucket should be set to private and URLs should be created via the `SignedURL` method to not expose the data.
    * `database` with `driver: inmemory|postgres`
        * `inmemory` is a simple non-persistant store. Only usable for quick debugging and in tests to mock basic database functionality.
    * `document-ai` `project-id`, `processor-id` and `location` should be set according to your processor endpoint. See: [DocumentAI Request](https://cloud.google.com/document-ai/docs/send-request#curl)
        * `credsfile` service account credentials. Make sure the service account has permissions for Document AI. See: [Google Cloud Service Accounts](https://developers.google.com/workspace/guides/create-credentials#service-account)

* Compile `make build` and run the executable `./expense-bot` or simply run the app by executing `go run cmd/expense-bot/main.go`
* A `Dockerfile` and `compose.yml` is included which include the dockerized version of the app and also the postgress database if needed.
    * **TODO** Add staged build and produce a light-weight image from the `scratch` image that is suited for production deplyoment.

## REST API
* The app implements a REST API for `expenses/`
* POST `expenses/?tags=tag1&tags=tag2...`
    * Payload `Content-Type: multipart/form-data` with a single `file` field
    * Add `tags` to a receipt with query parameters.
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
            "tags" : ["tag1", "tag2"]
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
            "tags" : ["tag1", "tag2"]
            "mime_type": "application/pdf",
            "path": "83bfe566-4254-4333-8ed1-7a54f918e796.pdf",
            "json_path": "83bfe566-4254-4333-8ed1-7a54f918e796.json"
        }
        ```
        * `id` generated `UUID`
        * `filename` original upload filename
        * `status` the receipt status `pending`, `ready`, `failed`
            * When a new receipt is uploaded an entry with `pending` status is created
            * The document processor sends a request to the processor API and updates with a `ready` status on completion
            * Or if the request timeouts or any other error is encountered the status is set to `failed`
        * `tags` list of tags associated with the receipt
        * `mime_type` uploaded file MIME Type. See supported formats: [DocumentAI file types](https://cloud.google.com/document-ai/docs/file-types)
        * `path` and `json_path` are the stored filenames. Does not include the `store.location` directory. But using the same file store it will retrieve the file correctly. This should be reworked in a complete app and is only a demo version.
* GET `expenses/?tags=tag1&tags=tag2`
    * Get receipts with any of the tags
    * Could add a paramater to get intersection or union - get only receipts with all the tags or get receipts with any of the tags.
    * **TODO:** will return correct Receipts but only with tags from query. Simple fix - see comments in `postgres.go` implementation.

## Expense Engine
The **Expense Engine** acts as a pipeline and listen&dispatch service. It asynchronously manages the processing of receipts uploaded via the **REST API**.
Currently the pipeline is hardcoded as processes within the pipeline have a linear progression from start to end with no way to alter and configure the pipeline and the order of execution. Some processes could very well be executed in parallel like **Translation** and **Currency Conversion** as they in no way rely on the result of eachother. The only change in order occurs if any of the steps fail and return an error - the pipeline will stop the process and mark the Receipt status as failed. This could be massively improved by identifying recoverable errors - receipt processor via Azure failed? Try the same with Google. Simply send a `EventMsg new` with instructions to use a particular processor to the **Expense Engine** and the process will start over.

### Pipeline
* `new` - a document was uploaded by the REST API and the process can start. It is sent to a receipt processor like Google Document AI or Azure Document Intelligence
* `processed` - the processor has finished and returned raw data. Dispatch data transformation to parse the data into a common **Expense** type
* `transformed` - the data is now transformed into a common data structure and post-processing can be applied. Translation and Currency Conversion. Both translation and currency conversion depend on the parsed data. If the processor was not able to detect the transaction time, currency used and totals, taxes and other money values then the currrency post-process is skipped entirely as it has no data to work with. A processor will still return a lot of valuable data even if it is not detected correctly and identified as a relevant field. A better data transformation layer can improve that and make resonable guesses about the raw data. Different post processing could be ideally done in parallel. Just need to ensure that the transforms are orthogonal - the do not share any field between them so the order of applying of the post-processing transforms should not alter the result.
* `done` - the last step of the pipeline has finished successfully as all before than and the receipt is fully processed.
* `failed` - any of the steps in the pipeline failed and was not recovered from. The receipt is marked as failed and no further processing is to be done.
    
## Improvements & Scalability
* While still only a simple Demo/Test app, I for the most part tried to make it as functional and clean as possible.
* Thread safety. The Gin framework and the use of `go routines` in the `GoogleDocumentAI.Process` method are almost guaranteed to cause panics related to concurrent reads/writes. For example the `inmemory` database implemented on `map[uiid]Receipt` will cause panics in when multiple request will be handled at the same time. It could be solved by using `mutex` or thread safe maps or other thread safe solutions.
* All services and clients are handled via an interface. Thus it should be relatively easy to swap out the FileSystem based file storage with a cloud based storage bucket or any other solution by simply creating an adapter with the appropriate wrapper that implements the interface.
* Document processing can take some time up to several seconds and maybe longer for larger documents. Currently the app creates a receipt of a request and initializes it with a `pending` status. While the user could refresh the API `expense/` endpoint with their receipt `UUID` it would be better that the initial request also adds a `callBackURI` that the client provides to recive a notification on completion or failure.
* In a real world application I believe it would be best to separate the upload and processing parts in separate micro services which would communicate via a message broker, pub/sub or any other method. The processor service could have a worker pool architecture and and subscribe to processing request messages. One could then spin up as many processors on say K8s to scale according to demand.
* It would also be a good practice to have separate data models for the `Receipt` type. One for internal database and storage and one for exposure to API endpoints. Keeping them separate adds more verbosity and some duplication of code but in a complex application having different types for the same data based on context often makes more sense than one-fits-all solution.
* I have no delusions of grandeur that my code will have no bugs, never crash or not have performance issues. So a production ready solution should have more robust logging, performance metric collection by `Grafana` and some issue/fault tracker with `Sentry` or similar solutions. To be able to monitor the performance and help fix issues.
    