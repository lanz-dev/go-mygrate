# go-mygrate

[![Go Reference](https://pkg.go.dev/badge/github.com/lanz-dev/go-mygrate.svg)](https://pkg.go.dev/github.com/lanz-dev/go-mygrate)
[![Coverage Status](https://coveralls.io/repos/github/lanz-dev/go-mygrate/badge.svg?branch=main)](https://coveralls.io/github/lanz-dev/go-mygrate?branch=main)
[![Github Action](https://github.com/lanz-dev/go-mygrate/actions/workflows/main.yml/badge.svg)](https://github.com/lanz-dev/go-mygrate/actions/workflows/main.yml)

## Mμgrate

Mμgrate is a migration tool. It has no external deps and uses go stdlib. A migration in Mμgrate is an ID and a pair of
functions which will be executed in the order as they were registered.

It's the right tool if you are looking for a simple, small and handy solution to migrate different "things" in a project
up and down. Most people will use it for databases, but you can theoretically migrate anything with it.

**Features**

- ship your migrations compiled in your binary
- migrate programmatically
- no dependencies on any ORM
- use the same deps and drivers you are already using in the project
- migrate whatever you want! A migration is just an ID and a pair of functions which getting called in same order!
- ships with a json based FileStore, database/sql stores (SQLite/MySQL tested) and a MemoryStore
- caution: default stores using a mutex locking. It's not safe to run migrations on multiple instances with the same
  database at once!
    - but it's really easy to implement your own store which implements your correct locking mechanics
- there is no magic involved!

### Installation

```bash
go get github.com/lanz-dev/go-mygrate
```

### Usage

See the example folder. Have fun and build something great!
