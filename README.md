# promise-sub
A golang webserver consumer implementation of the subscriber of promise-pub. The reason golang is used here is its known concurrent capabillities making it a great fit for the use case where I need to asyncronously receive thousands of messages through a message broker and save them to persistent data store.

This webserver listens on port `8080`

## Protocols
It makes more sense to define protocols for machines to communicate with each other. This consumer implements two protocols for its users:
  1. Live protocol: the messages are processed in real-time one by one.
  2. Batch protocol: the messages are processed in batch all together.

## Workflow
Designing a distributed system is hard. thus we need to define how two machines can communicate with each other. This system works as follows:
  1. The client requests to open a new session
  2. The consumer creates a new session with a random V4 UUID, opens a new rabbitmq connection and a new queue with its name identical to the session uuid.
  3. The consumer returns the newly created session object to the client.
  4. If the client wishes to work in live mode, then it tells the consumer to start listening. If otherwise the client does nothing yet.
  5. The client then starts sending data over the queue
  6. If the client wishes to work in batch mode, then it starts telling the consumer to consume the queue
  7. The consumer consumes the queue and saves them to the database.
  8. The client terminates the session.

## Instalation
  1. clone this repo
  2. set the `GOPATH` to be inside this repo for example `export GOPATH=/path/to/repo/`
  3. set the `GOBIN` to the bin folder for example `export GOBIN=$GOPATH/bin/`
  4. run `go get` to install external libraries
  5. run `go run main.go`

## Modules
  - main: The entry point of this webserver
  - session: Manages the creation, deletion and processing of message queues
  - rabbit: Manages the low-level communication with rabbitmq

## Dependencies
  1. github.com/streadway/amqp for rabbitmq connections
  2. github.com/lib/pq for postgres connections
  3. github.com/nu7hatch/gouuid for generating UUID
