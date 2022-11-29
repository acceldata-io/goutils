# sysd

example app to demo libsysd capabilities

---

Compile `sysd` binary for Linux:

```shell
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -a -installsuffix cgo '-ldflags=-w -s' -o sysd .
```
