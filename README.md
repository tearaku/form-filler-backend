# Backend for Form Filler frontend
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
