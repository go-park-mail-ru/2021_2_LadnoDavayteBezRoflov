[![CI](https://github.com/go-park-mail-ru/2021_2_LadnoDavayteBezRoflov/actions/workflows/CI.yml/badge.svg?branch=main)](https://github.com/go-park-mail-ru/2021_2_LadnoDavayteBezRoflov/actions/workflows/CI.yml)

[![Go Report Card](https://goreportcard.com/badge/github.com/go-park-mail-ru/2021_2_LadnoDavayteBezRoflov)](https://goreportcard.com/report/github.com/go-park-mail-ru/2021_2_LadnoDavayteBezRoflov)

# Trello

Trello backend repository for LadnoDavayteBezRoflov team, autumn of 2021.

### Team

* [Anton Chumakov](https://github.com/TonyBlock);
* [Alexander Orletskiy](https://github.com/Trollbump);
* [Georgiy Sedoykin](https://github.com/GeorgiyX);
* [Dmitriy Peshkov](https://github.com/DPeshkoff).

### Mentors

* [Ekaterina Alekseeva](https://github.com/yletamitlu) — frontend mentor;
* [Roman Gavrilenco](https://github.com/gavroman) — frontend mentor;
* [Timofey Razumov](https://github.com/TimRazumov) — backend mentor.

### Frontend repository
[Link to frontend repository](https://github.com/frontend-park-mail-ru/2021_2_LadnoDavayteBezRoflov).

### API
[Link to API](https://app.swaggerhub.com/apis/DPeshkoff/LadnoDavayteBezRoflov).

### Deploy
[Link to deploy](https://brrrello.ru).

### Usage

> Starting the bare server (requires Redis and PostgreSQL): `go build ./cmd/api && sudo ./api`

> Starting the server using docker-compose: `docker-compose up`

> Running backend tests: `go test --coverpkg=$$(go list ./... | xargs echo | tr ' ' ,) ./... && go tool cover -func=.coverprofile`

### Directory structure

```bash
2021_2_LadnoDavayteBezRoflov
|--cmd/api
|  |-main.go
|  |-server.go
|  |-settings.go
|
|--app
|  |-handlers
|  |-models
|  |-repositories
|  |-usecases
|
|--pkg
```
