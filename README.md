# GoPlug

[![Build Status](https://travis-ci.org/yurigorokhov/goplug.svg)](https://travis-ci.org/yurigorokhov/goplug)


HTTP Plug library using net/http. Simplifies accessing of web resources over HTTP.

- No Dependencies
- Fluent, mutable interface

```go

import (
	plug "github.com/yurigorokhov/goplug"
)

p, err := plug.New('http://google.com')
r := <-p.With('q', 'my search').With('limit', '10').Get()

```
