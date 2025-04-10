# go-datatrails-common

Public repository for base Go utility packages.

> [!WARNING]
> No guarantee is given that packages in this repo will not be removed.
> Import this repo at your own risk.

# Multirepo development

Use go work spaces for co-development against another repo that may import this

```bash
cd
mkdir workspace
cd workspace
git clone https://github.com/datatrails/go-datatrails-common.git
git clone https://github.com/datatrails/go-datatrails-logverification.git
go work init
go work use go-datatrails-common
go work use go-datatrails-logverification
```

