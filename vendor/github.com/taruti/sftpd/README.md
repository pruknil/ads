# sftpd - SFTP server library in Go

# License: MIT - Docs [![GoDoc](https://godoc.org/github.com/taruti/sftpd?status.png)](http://godoc.org/github.com/taruti/sftpd)

# Recent changes - 2019
+ `Attr.FillFrom` cannot fail and now does not return an error value. Previously it was always nil.

# FAQ

## ssh: no common algorithms

The client and the server cannot agree on algorithms. Typically this
is caused by an ECDSA host key. If using sshutil add the
``sshutil.RSA2048`` flag.

## Enabling debugging output

```
go build -tags debug -a
```

Will enable debugging output using package `log`.

# TODO
+ Renames
+ Symlink creation
