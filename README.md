# Identity API
Handles the user identity, authentication, authorization.

This is still very much WIP. For now it can authenticate existing users and give them a refresh token.

## Setting up for development
You need at least golang 1.15 https://golang.org/.

Install dependencies with `go get ./...` 

To run the app for development:
`JWT_TOKEN=secret go run main.go`

## Building docker container
`docker build --tag "identity-api:latest" .`

## Starting with database
`JWT_TOKEN=secret docker-compse up`

