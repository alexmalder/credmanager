# Credmanager

## How to build

```bash
git clone https://github.com/alexmalder/credmanager
cd credmanager
go install
go build -o credmanager main.go
```

## Example usage

```bash
$ go run main.go --help

Usage: main [-c <config>]

global description

Options:
    -c, --config=<config>   key of a new secret (default: config.yml)
    -h, --help              usage (-h) / detailed help text (--help)

Available commands:
    create-file             create key-value pair as file
    create-value            create key-value pair as string
    delete                  delete secret by key
    drop                    drop secrets table
    get                     get secret by key
    put-file                put secret by key
    put-value               put secret by key
    select                  select secrets
```

## Environment variables

### Bitwarden intergration

- `BITWARDEN_BACKUP_PATH`: path for bitwarden json backup file

### GPG

- `GPG_HOMEDIR`: path for a gpg directories
- `GPG_PASSPHRASE`: gpg passphrase

### Postgres

- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `POSTGRES_HOST`
- `POSTGRES_PORT`
- `POSTGRES_DB`


## Todo

- [ ] check put
- [ ] check delete
