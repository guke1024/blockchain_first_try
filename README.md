### This project has only been tested under Windows operating system and go 1.17

## Pull project
```shell
git clone https://github.com/guke1024/blockchain_first_try.git
```

## Add required modules and build
```shell
cd blockchain_first_try
rm .\blockChain.db
go get go.etcd.io/bbolt/...
go build blockchain_first_try
```
## Run
```shell
.\blockchain_first_try.exe
```

## Usage
```shell
.\run.bat [command]
```