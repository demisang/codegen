# codegen

Code generator

### Install

```shell
go get github.com/demisang/codegen
cp $GOPATH/pkg/mod/github.com/demisang/codegen/cmd/codegen-service/main.go cmd/codegen-service/main.go
cp $GOPATH/pkg/mod/github.com/demisang/codegen/templates-example internal/codegen/templates-example
```

Makefile:

```makefile
## codegen: run codegen app on http://localhost:4765
codegen:
    go run cmd/codegen-service/main.go --root="./" --templates="./internal/codegen/templates"
```

Open GUI http://localhost:4765

### Usage

Feel-free to create/edit templates in `internal/codegen/templates` adapt for your project
