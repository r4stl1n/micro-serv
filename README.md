# Micro Serv

Simple microservice framework that relies on nats and offers a cache and data storage service. This framework is designed to quickly get up and running with basic microservices based on nats micro.

## Installation

Install GO version 1.21.x or greater

```bash
go mod tidy
```

## Docker Development Usage

To recreate every docker image and redeploy
```bash
./scripts/deployment/debug-local.sh all
```

To recreate and deploy a specific image
```bash
./scripts/deployment/debug-local.sh micro-data-service
```

To deploy and have all services in the foreground
```bash
./scripts/deployment/debug-local.sh
```

## Deployment

For a quick and drity deployment using docker-compose
```bash
./scripts/deployment/deploy-local.sh
```

## Test
Due to the tests spinning up temporary services, they expect the core services to be online and running.

```bash
./scripts/deployment/start-core-services.sh
```

Run test cases example
```bash
MONGO_HOST=localhost:27017 MONGO_PASS=toor MONGO_USER=root NATS_HOST=localhost:4222 NATS_PASS=T0pS3cr3t NATS_USER=ruser REDIS_HOST=localhost:6379 REDIS_PASS=toor go test ./...
```

## Deployment
Each service relies on enviornment variables in order to operate and identify the required services.

### Dummy-Service
* NATS_HOST=<host>
* NATS_PASS=<pass>
* NATS_USER=<user>

### Data-Service
* NATS_HOST=<host>
* NATS_USER=<user>
* NATS_PASS=<pass>
* MONGO_HOST=<host>
* MONGO_USER=<user>
* MONGO_PASS=<pass>

### Cache-Service
* NATS_HOST=<host>
* NATS_USER=<user>
* NATS_PASS=<pass>
* REDIS_HOST=<host>
* REDIS_PASS=<pass>

## Note
This project uses a special image of nats created to support 32mb nats messages. Nats supports up to 64mb messages however it is not recommended to go that far. That image is r4stl1n/nats-large-payload:latest

## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)