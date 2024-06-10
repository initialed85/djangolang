# djangolang

Big WIP.

```shell
# shell 1
./run-env.sh

# shell 2
PORT=7070 POSTGRES_DB=some_db POSTGRES_PASSWORD=some-password go run ./cmd/main.go serve

# shell 3
websocat ws://localhost:7070/__stream | jq

# shell 4
find . -type f -name '*.*' | grep -v '/model_generated/' | entr -n -r -cc -s "DJANGOLANG_DEBUG=1 POSTGRES_DB=some_db POSTGRES_PASSWORD=some-password go test -v -failfast -count=1 ./pkg/template && DJANGOLANG_DEBUG=1 POSTGRES_DB=some_db POSTGRES_PASSWORD=some-password go test -v -failfast -count=1 ./pkg/model_generated_test"
```
