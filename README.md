# [Multi-level logger based on go std log](https://golangexample.com/multi-level-logger-based-on-go-std-log/)

[![License](https://img.shields.io/badge/license-BSD-green)](https://github.com/ccpaging/log/blob/master/LICENSE) [![GoDoc](https://godoc.org/github.com/ccpaging/log?status.svg)](https://godoc.org/github.com/ccpaging/log) [![Build Status](https://github.com/ccpaging/log/actions/workflows/go.yml/badge.svg)](https://github.com/ccpaging/log/actions/workflows/go.yml) [![Maintainability](https://codeclimate.com/github/ccpaging/log/badges/gpa.svg)](https://codeclimate.com/github/ccpaging/log/maintainability)

## mlog

the mlog is multi-level logger based on go std log. It is:

* Simple
* Easy to use
* Easy to extend a new logger like color log, syslog, rolling file log ...

NOTHING ELSE

```
package main

import (
    log "github.com/ccpaging/log/mlog"
)

func main() {
    log.Debug("This is Debug")
    log.Info("This is Info")
}
```

Issues:

[#48464](https://github.com/golang/go/issues/48464) [#48503](https://github.com/golang/go/issues/48503) [#28412](https://github.com/golang/go/issues/28412)  [#13182](https://github.com/golang/go/issues/13182) [#28327](https://github.com/golang/go/issues/28327) [#32062](https://github.com/golang/go/issues/32062)
