homedir=./
gpg --gen-key --homedir $homedir
gpg --no-default-keyring --homedir $homedir/ --export-secret-keys > $homedir/secring.gpg
gpg --no-default-keyring --homedir $homedir/ --export > $homedir/pubring.gpg
