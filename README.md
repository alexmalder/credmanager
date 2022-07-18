# Credmanager

## How to build

```bash
git clone https://github.com/alexmalder/credmanager
cd credmanager
go install
go build -o credmanager main.go
```

## Hot to generate keys

```bash
export GPG_HOMEDIR=gpg # directory in gitignore
gpg --expert --full-gen-key --homedir $GPG_HOMEDIR
gpg --no-default-keyring --homedir $GPG_HOMEDIR/ --export-secret-keys > $GPG_HOMEDIR/secring.gpg
gpg --no-default-keyring --homedir $GPG_HOMEDIR/ --export > $GPG_HOMEDIR/pubring.gpg
```

## Example usage

```bash
go run main.go --help
Usage: main [-c <config>]

global description

Options:
    -c, --config=<config>   key of a new secret (default: config.yml)
    -h, --help              usage (-h) / detailed help text (--help)

Available commands:
    create-file             create key-value secret as file
    create-value            create key-value secret as value
    drop                    drop secrets table
    get                     get secret by key
    import-bitwarden        import bitwarden json file
    migrate                 create tables if does not exists
    put-file                put secret by key
    put-value               put secret by key
    select                  select secrets
```

## Database table structure

| Field      | Type          | Note                             |
| ---------- | ------------- | -------------------------------- |
| key        | VARCHAR(255)  | key of a secret                  | 
| revision   | INTEGER       | number of revision               | 
| value      | VARCHAR(4096) | value of a secret                | 
| username   | VARCHAR(255)  | optional field of a username     |
| uri        | VARCHAR(1024) | optional field of a uri          |
| notes      | VARCHAR(4096) | optional field of a notes        |
| type       | VARCHAR(8)    | "file", "login" or custom        |
| is_deleted | BOOLEAN       | set `is_deleted` for hiding      |


## Environment variables

### GPG

- `GPG_HOMEDIR`: path for a gpg directories
- `GPG_PASSPHRASE`: gpg passphrase
- `GPG_SECRING`: gpg secring
- `GPG_PUBRING`: gpg pubring

### Postgres

- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `POSTGRES_HOST`
- `POSTGRES_PORT`
- `POSTGRES_DB`

## TODO

- [ ] update keyring
- [ ] restore version

## Authors

- `vnmntn@mail.ru`
