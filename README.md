# Golang project boilerplate
This repo's name explain itself. Still have a lot of things todo

## TODOs
- [ ] Add license file
- [ ] Write document about how to use `.env` file
- [ ] Makefile to run, build

and more

## Migration commands
Create migration in `migrations/sql` directory
```
go run cmd/migrator/migrator.go create -n my_file_name
```
Execute all migrations
```
go run cmd/migrator/migrator.go up
```
