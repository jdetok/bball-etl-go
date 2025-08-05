# bball-etl-go
postgres db build & extract-transform-load process written in Go to insert nba/wnba statistics from [nba.com](stats.nba.com) into database. the database is connected to the backend of [jdeko.me/bball](https://jdeko.me/bball), a basketball stats site i built along with my personal website
### jdeko.me repo: [jdetok/go-api-jdeko.me](https://github.com/jdetok/go-api-jdeko.me)

## /etl: go extract-transform-load
## /sql: database build scripts
## [build documentation](z_docs/BUILD.md)
- bld.sh builds/runs postgres docker container, runs go etl program, & calls stored procedures to build & insert data into database from one script
## [database design/architecture documentation](z_docs/database/schemas.md)