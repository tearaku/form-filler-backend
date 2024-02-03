# Backend for Form Filler frontend
## Running Tests
For db tests, stub data must be inserted into an existing postgres db first.

Run the test with the following env variable set to generate stub data (note: do NOT do this against an actual db!).
```bash
GEN_TEST_DATA=1 go test -v ./...
```

### `godotenv` is erroring on tests
You may need to change `repoDir` in `internal/dataSrc/helper.go` file.

See the linked issue in said file for details on why this is a thing.

## Install & Run
```bash
docker-compose up
```
- Provide the necessary environment variables listed in the compose file (as .env file)
- Get the Chinese font from [here](https://data.gov.tw/dataset/5961), and install the `Kai` family into the `resources` directory (otherwise outputted PDF will not render Chinese properly)

# To-do List:
- [] Universal API-safe error handling w/ semantics
- [] Add support for absurd team size (30+)
    - Right now each sheet has an implicit limit on existing rows --> if sizes go over this, there will be issues
- [] Add support for gRPC
