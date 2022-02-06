### This project has only been tested with go 1.17 and golang on Windows operating system.

# Usage

## Pull project
```shell
git clone https://github.com/guke1024/blockchain_first_try.git
```

## Add required modules and build
First, you need to install [btcd](https://github.com/btcsuite/btcd).
```shell
cd $GOPATH\pkg\mod\github.com
git clone https://github.com/btcsuite/btcd.git
go env -w GO111MODULE=on
go install -v . ./cmd/...
```
Second, you need to install [bbolt](https://github.com/etcd-io/bbolt).
```shell
go get go.etcd.io/bbolt/...
```
Finally, you need to initialize go mod.
```shell
cd blockchain_first_try
go mod init
go mod tidy
.\run.bat
```
## Run
```shell
.\blockchain_first_try.exe [command]
```

### Or you can use the following command to test quickly.
```shell
.\blockchain_first_try.exe transfer 1JLyj9EK6sEv2NpVFV74NvEmgakc2uqmUj 134cxCG8v3RXjfzQRfrQbfw9kKc8bcwh74 20 16f7nqgn2nxnQkcdGtfSBCr1rparBDMPdH wuhu~
.\blockchain_first_try.exe getBalance --address 1JLyj9EK6sEv2NpVFV74NvEmgakc2uqmUj
.\blockchain_first_try.exe getBalance --address 134cxCG8v3RXjfzQRfrQbfw9kKc8bcwh74
.\blockchain_first_try.exe getBalance --address 16f7nqgn2nxnQkcdGtfSBCr1rparBDMPdH
```
If you don't want to use this data, please delete the `wallet.dat`.
