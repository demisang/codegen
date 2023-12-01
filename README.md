# codegen

Code generator

![codegen_screenshot](https://github.com/demisang/codegen/assets/600251/3657de47-72d8-4027-8dc8-2c88f181806e)

### Install


```shell
go install github.com/demisang/codegen@latest
```

### Get template samples

* Golang onion-architecture code samples (will extract to `./pkg/codegen/templates`)
  ```shell
  DIR=pkg/codegen/templates; mkdir -p $DIR && wget -qO- https://github.com/demisang/codegen/raw/main/templates/go_onion.tar.gz | tar xvzf - -C $DIR
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
