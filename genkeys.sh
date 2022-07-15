#!/bin/bash
gpg --expert --full-gen-key --homedir $GPG_HOMEDIR
gpg --no-default-keyring --homedir $GPG_HOMEDIR/ --export-secret-keys > $GPG_HOMEDIR/secring.gpg
gpg --no-default-keyring --homedir $GPG_HOMEDIR/ --export > $GPG_HOMEDIR/pubring.gpg
