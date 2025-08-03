# TODO: sunday 8/3 5pm
- WRITE TESTS
    - validate headers from response match intake column names
    - validate seasons
- decide how api integration should work
    - queries or a table (like in mariadb)
    - try query first
        - BUILD THE QUERY
- figure out how to separate (and have accessible) the version of main that runs for all seasons & the version that will run nightly once up & running
    - probably will need to make the etl process into a package then different packages will call it in different ways
