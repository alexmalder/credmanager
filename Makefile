test:
	go build -o credmanager main.go
	#./credmanager create-file -k clickhouse_config -f test.yml
	#./credmanager get -k clickhouse_config
	./credmanager create-value -k postgres_user -v postgres_password --username postgres --uri "127.0.0.1" --notes "my-local-database"
	./credmanager get -k postgres_user
	./credmanager select
	./credmanager put-value -k postgres_user
	./credmanager put-value -k postgres_user --notes my-local-db
	./credmanager get -k postgres_user
	./credmanager put-value -k postgres_user --uri 127.0.0.1
	./credmanager put-value -k postgres_user --uri 127.0.0.2
	./credmanager get -k postgres_user
	./credmanager put-value -k postgres_user --value postgres_password
	./credmanager put-value -k postgres_user --value postgres_pswd
	./credmanager get -k postgres_user
	./credmanager put-value -k postgres_user --username postgres
	./credmanager put-value -k postgres_user --username postgres_u
	./credmanager get -k postgres_user
	./credmanager put-value -k postgres_user --is_deleted false
	./credmanager put-value -k postgres_user --is_deleted true
	./credmanager get -k postgres_user

migrate:
	go run main.go migrate

drop:
	go run main.go drop

install:
	go build -o credmanager main.go
	sudo mv ./credmanager /usr/local/bin
