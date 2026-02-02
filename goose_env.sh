export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
GOOSE_DRIVER=postgres GOOSE_DBSTRING=postgres://program:Password123@127.0.0.1:5432/chatdb GOOSE_MIGRATION_DIR=./migrations GOOSE_TABLE=public.goose_migrations "$@"