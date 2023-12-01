# codegen

Code generator

### Install

```shell
go install github.com/demisang/codegen@latest
```

### Get template samples
```shell
# Golang onion-architecture code samples (will extract to ./pkg/codegen/templates)
wget -qO- https://github.com/demisang/codegen/tree/main/templates/go_onion.tar.gz | tar xvz - -C pkg/codegen/templates
```

Makefile:

```makefile
## codegen: run codegen app on http://localhost:4765
codegen:
    codegen --root="./" --templates="./pkg/codegen/templates"
```

Open GUI http://localhost:4765

### Usage

Feel free to create/edit templates in `pkg/codegen/templates` adapted for your project
