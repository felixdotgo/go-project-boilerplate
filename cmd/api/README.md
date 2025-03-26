# api service
## Run service
Under root directory, run the following command
```
make run CMD=api
```

## Migration commands
Create migration in `migrations/sql` directory
```
go run cmd/migrator/migrator.go create -n my_file_name
```
Execute all migrations
```
go run cmd/migrator/migrator.go up
```
