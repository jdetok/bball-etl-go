# bball-etl-go
extract-transform-load process written in Go to insert nba/wnba statistics from [nba.com](stats.nba.com) into a postgres database. the database is connected to the backend of [jdeko.me/bball](https://jdeko.me/bball), a basketball stats site i built along with my personal website
### database repo: [jdetok/postgres-bball-db](https://github.com/jdetok/postgres-bball-db)
### jdeko.me repo: [jdetok/go-api-jdeko.me](https://github.com/jdetok/go-api-jdeko.me)

# extrct package
- http requests to stats.nba.com

# trsfrm package
- unmarshal reutrned json into structs

# pgload package
- process & load the extracted data into the postgres database