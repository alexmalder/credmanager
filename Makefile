test:
	go build -o credmanager main.go
	./credmanager create-value -k postgres_user -v postgres_password
	./credmanager create-file -k clickhouse_config -f test.yml
	./credmanager get -k postgres_user
	./credmanager get -k clickhouse_config
	./credmanager select
	#./credmanager delete -k postgres_user
	#./credmanager delete -k clickhouse_config
	#./credmanager select
