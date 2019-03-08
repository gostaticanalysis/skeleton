# skeleton 

skeleton is create skeleton codes for golang.org/x/tools/go/analysis.

## Insall

```
$ go get -u github.com/gostaticanalysis/skeleton
```

## How to use

### Create skeleton codes in GOPATH

```
$ skeleton pkgname
pkgname
├── cmd
│   └── pkgname
│       └── main.go
├── pkgname.go
├── pkgname_test.go
└── testdata
    └── src
        └── a
            └── a.go
```

### Create skeleton codes with import path

```
$ skeleton -path="github.com/gostaticanalysis/pkgname"
pkgname
├── pkgname.go
├── pkgname_test.go
└── testdata
    └── src
        └── a
            └── a.go
```

### Create skeleton codes without cmd directory

```
$ skeleton -cmd=false pkgname
pkgname
├── pkgname.go
├── pkgname_test.go
└── testdata
    └── src
        └── a
            └── a.go
```
