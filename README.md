# kaaskill

## install
```sh
git clone https://github.com/kotajacob/kaaskiller && cd kaaskiller
go build
```

## usage
```sh
export CIVO_API_KEY="somefancyapikey"
./kaaskiller
```

It will fetch a list of all active Kubernetes clusters for your account and
allow you to enter a number to delete one of them.
