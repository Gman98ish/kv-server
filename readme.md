Simple KV server that safely handles parallel reads and writes

## Building and running

```
$ go build
$ export PORT=${your port here}
$ ./key-value-server
```

## Testing

Can use `go test` or `go test -race` if you want to also test for race conditions