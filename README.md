# CoreStore packages | WIP (Work in Progress)

[![Build Status](https://travis-ci.org/corestoreio/csfw.svg?branch=master)](https://travis-ci.org/corestoreio/csfw) [![wercker status](https://app.wercker.com/status/d7d0bdda415d2228b6fb5bb01681b5c4/s/master "wercker status")](https://app.wercker.com/project/bykey/d7d0bdda415d2228b6fb5bb01681b5c4) [![Appveyor Build status](https://ci.appveyor.com/api/projects/status/lrlnbpcjdy585mg1/branch/master?svg=true)](https://ci.appveyor.com/project/SchumacherFM/cspkg/branch/master) [![GoDoc](http://godoc.org/github.com/corestoreio/csfw?status.svg)](http://godoc.org/github.com/corestoreio/csfw) [![Join the chat at https://gitter.im/corestoreio/csfw](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/corestoreio/csfw?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge) [goreportcard](http://goreportcard.com/report/Corestoreio/csfw)

eCommerce packages which is NOT compatible to a Magento 2 database schema but kinda.

Magento is a trademark of [MAGENTO, INC.](http://www.magentocommerce.com/license/).

Min. Go Version: 1.9

## Usage

To properly use the CoreStore packages, some environment variables must be set
before running `go generate`. (TODO)

### Required settings

`CS_DSN` the environment variable for the MySQL connection.

```shell
$ export CS_DSN='magento2:magento2@tcp(localhost:3306)/magento2'
```

```
$ go get github.com/corestoreio/csfw
$ export CS_DSN='see previous'
$ cd $GOPATH/src/github.com/corestoreio/csfw
$ xgo run xcodegen/tableToStruct/*.go
```

## Testing

Setup database for Magento 2 and fill them with
the provided [testdata](https://github.com/corestoreio/cspkg/tree/master/testData).

Create a DSN env var `CS_DSN` and point it to Magento 2 database. Run the tests. 

## TODO

If you find an entry in the source like `TODO(CS)` that means you can ask `CS`
to get more information about how to implement and what to fix if the context of
the todo description isn't understandable.

- Create Magento 2 modules to setup test database and test Magento system.

## Details

* **single repo**. CoreStore is a single repo. That means things can be
    changed and rearranged globally atomically with ease and
    confidence.

* **no backwards compatibility**. CoreStore makes no backwards compatibility
    promises. If you want to use CoreStore, vendor it. And next time you
    update your vendor tree, update to the latest API if things in CoreStore
    changed. The plan is to eventually provide tools to make this
    easier.

* **forward progress** because we have no backwards compatibility,
    it's always okay to change things to make things better. That also
    means the bar for contributions is lower. We don't have to get the
    API 100% correct in the first commit.

* **no Go version policy** CoreStore packages are usually built and tested
    with the latest Go stable version. However, CoreStore has no overarching
    version policy; each package can declare its own set of supported
    Go versions.

* **code review** contributions must be code-reviewed.

* **CLA compliant** contributors must agree to the CLA.

* **docs, tests, portability** all code should be documented in the
    normal Go style, have tests, and be portable to different
    operating systems and architectures. We'll try to get builders in
    place to help run the tests on different OS/arches. For now we
    have Travis at least.

## Contributing

Please have a look at the [contribution guidelines](https://github.com/corestoreio/corestore/blob/master/CONTRIBUTING.md) That document is an idea.

## Acknowledgements

Some packages have been fully refactored but the initial idea has been extracted
from the geniuses listed below:

| Name | Package | License |
| -------|----------|-------|
| Steve Francia | [util/conv](http://github.com/corestoreio/cspkg/tree/master/utils/conv) | MIT Copyright (c) 2014 |
| Martin Angers and Contributors. | [ctxthrottled](http://github.com/corestoreio/cspkg/tree/master/net/ctxthrottled) | The MIT License (MIT) 2014 |
| Dave Cheney <dave AT cheney.net> | [util/errors](https://github.com/pkg/errors) | The MIT License (MIT) 2015 |
| Jad Dittmar | [finance](https://github.com/Confunctionist/finance) aka. [money](http://github.com/corestoreio/cspkg/tree/master/storage/money) | Copyright (c) 2011 |
| Wenbin Xiao | [util/sqlparser](https://github.com/xwb1989/sqlparser) | Copyright 2015 BSD Style |
| Google Inc | [youtube/vitess\sqlparser](https://github.com/youtube/vitess) | Copyright 2012 BSD Style |
| Olivier Poitrey| [ctxmw.WithAccessLog](https://github.com/corestoreio/cspkg/tree/master/net/ctxmw) & CORS | Copyright (c) 2014-2015  MIT License |
| Dave Grijalva| [csjwt](https://github.com/corestoreio/cspkg/tree/master/util/csjwt) | Copyright (c) 2012 MIT License |
| Uber Technologies, Inc. | [log](https://github.com/corestoreio/cspkg/tree/master/log) | Copyright (c) 2016 MIT License |
| 2013 The Go Authors | [singleflight](https://github.com/corestoreio/cspkg/tree/master/sync/singleflight) | Copyright (c) 2013 BSD Style |
| Ventu.io, Oleg Sklyar, contributors. | [util/shortid](http://github.com/corestoreio/cspkg/tree/master/utils/shortid) | MIT License Copyright (c) 2016, |
| Carl Jackson (carl@avtok.com) (Goji) | [net/responseproxy](http://github.com/corestoreio/cspkg/tree/master/net/responseproxy) | Copyright (c) 2014, 2015, 2016 |
| Greg Roseberry, 2014; Patrick O'Brien, 2016 | [util/null](http://github.com/corestoreio/cspkg/tree/master/util/null) | BSD Copyright (c) 2014, 2015, 2016 |
| The Go-MySQL-Driver Authors | [util/null/time_mysql.go](http://github.com/corestoreio/cspkg/tree/master/util/null/time_mysql.go) | Mozilla Public License, v. 2.0, Copyright 2012  |
| siddontang | [storage/binlogsync](http://github.com/corestoreio/cspkg/tree/master/storage/binlogsync) | MIT Copyright (c) 2014  |
| siddontang | [storage/myreplicator](http://github.com/corestoreio/cspkg/tree/master/storage/myreplicator) | MIT Copyright (c) 2014  |
| Tace De Wolf | [util/byteconv](http://github.com/corestoreio/cspkg/tree/master/util/byteconv) | MIT Copyright (c) 2015  |
| Copyright 2013 The Camlistore Authors | [util/byteconv](http://github.com/corestoreio/cspkg/tree/master/util/byteconv) | Apache 2.0  |

## Licensing

CoreStore is licensed under the Apache License, Version 2.0. See
[LICENSE](https://github.com/corestoreio/corestore/blob/master/LICENSE) for the full license text.

## Copyright

[Cyrill Schumacher](https://cyrillschumacher.com) - [PGP Key](https://keybase.io/cyrill)
