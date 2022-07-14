#!/bin/bash
encrypt_salt() {
    tar --exclude /srv/salt/bin -zczf salt.tar.gz /srv/salt
    gpg --yes --encrypt --recipient $RECIPIENT salt.tar.gz
    rm salt.tar.gz
}

encrypt_salt
