# kaaskill

## install
```sh
git clone https://github.com/kotajacob/kaaskill && cd kaaskill
go build
```

## usage
```sh
export CIVO_API_KEY="somefancyapikey"
export DIGITAL_OCEAN_API_KEY="somefancyapikey"
export LINODE_API_KEY="somefancyapikey"
./kaaskill
```

The program first allows you to select which of the providers you are
interacting with. You need to have set API keys set for the providers you're
selecting.

It will then fetch a list of all active Kubernetes clusters for your account
with your specified provider and allow you to enter a number to delete one of
them.
