# Authentication service
## Run service
Under root directory, run the following command
```
make run CMD=svc-auth
```

## Migration commands
From the repo root:
```
go run ./cmd/svc-auth/cmd/migrator create -n my_file_name
```

From `cmd/svc-auth/`:
```
go run ./cmd/migrator/migrator.go create -n my_file_name
go run ./cmd/migrator/migrator.go up
```
