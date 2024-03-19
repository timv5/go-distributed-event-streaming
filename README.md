# go-distributed-event-streaming
Distributed event streaming written in Golang using RabbitMQ

## Description
The purpose of this project is learning golang programming language.

First call exposed endpoint. Message is saved to database and published to
RMQ queue. As soon as **producer** gets a message there is a goroutine in which rest call is made on external system.
As it is made in a separate goroutine, we dont wait for a response, but we start saving a message (& message history).
Consumer then receives this message and updates it status in db.
Project also contains **outbox patter**. Before message is sent from producer to consumer it's first saved to outbox
db table, so that we ensure transaction is completed. In main class goroutine picks up outbox records from db table
every 5s, unmarshall them and sends them on queue for consumer.

## Used technologies and prerequisites
- golang
- postgres
- docker, docker compose
- RabbitMQ

## How to run a project
- go to /docker
- run command: docker-compose -f docker-compose.yml up -d
- generate initial schema: generate_schema.sql
- run go app producer
- run go app consumer

## How to run a project - alternative
- go to branch /docker
- go to folder /docker
- run docker-compose -f docker-compose.yml up -d
- 4 containers will run
- application & db & rmq will be accessible from local machine as before

## RabbitMQ, Postgres Endpoints and functionality
- rabbitMQ accessible on: http://localhost:15672/#/ (user: guest, pass: guest)
- endpoints: http://localhost:8080/api/message/send
  - body: {"header": "he","body": "bo"}
