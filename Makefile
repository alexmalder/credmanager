test:
	go build -o credmanager main.go
	#./credmanager create-file -k clickhouse_config -f test.yml
	#./credmanager get -k clickhouse_config
	./credmanager create-value -k postgres_user -v postgres_password --username postgres --uri "127.0.0.1" --notes "my-local-database"
	./credmanager get -k postgres_user
	./credmanager select
	./credmanager put-value -k postgres_user
	./credmanager put-value -k postgres_user --notes my-local-db
	./credmanager put-value -k postgres_user --uri 127.0.0.1
	./credmanager put-value -k postgres_user --uri 127.0.0.2
	./credmanager put-value -k postgres_user --value postgres_password
	./credmanager put-value -k postgres_user --value postgres_pswd
	./credmanager put-value -k postgres_user --username postgres
	./credmanager put-value -k postgres_user --username postgres_u

drop:
	go run main.go drop
