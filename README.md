# `codex`
[![Go Reference](https://pkg.go.dev/badge/pkg.coulon.dev/codex.svg)](https://pkg.go.dev/pkg.coulon.dev/codex)
SSH certificates distribution infrastructure

**Warning:** codex is still under development

## Scope

Managing SSH permissions on multiple servers through the `authorized_keys` file can become tedious when you have a large set of users.
A common solution is to use OpenSSH certificates, which lets you authorize users without having to contact the targeted servers. However, manually assigning certificates can become tiresome, especially if you define short expiration times.

The goal of this project is to simplify the management of users while automating the distribution of certificates.

*If you want to know more about OpenSSH certificates, I recomend starting with [this article](https://engineering.fb.com/2016/09/12/security/scalable-and-secure-access-with-ssh/) by Meta engineers.*

## Architecture

codex implements a client-server infrastructure where the server stores a list of users and serves certificates to the clients.

User data is stored on disk using [bbolt](https://pkg.go.dev/go.etcd.io/bbolt), an embedded key-value store. User are represented as an SSH public key (their primary key) and the attached informations are their expiration, principals and the serial number of the last delivered certificate. User data can be managed from a command line utility named `codexctl`.

The codex server which is named `codexd` will listen for SSH connections on port `3646` by default. Clients will authenticate through their public key, the server will look for the corresponding user in its database and will return a new and signed OpenSSH certificate.

This procedure can either be done through the standard ssh client as shown below:
```shell
$ ssh codex.example.org -p 3646

  ___  _____  ____  ____  _  _
 / __)(  _  )(  _ \( ___)( \/ )
( (__  )(_)(  )(_) ))__)  )  (
 \___)(_____)(____/(____)(_/\_)
                         v0.0.0
Hello user,
Here is your new certificate:

ssh-ed25519-cert-v01@openssh.com [...]

Certificate details:

    Serial number: 42
    Expiration date: 08 Sep 22 11:43 CEST

Connection to localhost closed.
```
or using the `codex` utility which will automate the process of saving the new certificate in the user's keychain.
```shell
$ codex refresh foobar
Sucessfully updated your certificate at ~/.ssh/foobar-cert.pub
```
