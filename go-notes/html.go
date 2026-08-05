package main

import _ "embed"

// The file lives at static/index.html — edit it like any normal HTML file.
// At compile time, go:embed bakes it into the binary. No runtime file dependency.

//go:embed static/index.html
var indexHTML string
