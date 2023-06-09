# go-distributed-event-streaming
Distributed event streaming written in Golang using RabbitMQ

## Description
First call exposed endpoint. Message is saved to database and published to
RMQ queue. Consumer then receives this message and updates it status in db.

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
