package main

import (
	"github.com/onosproject/helmit/pkg/registry"
	"github.com/onosproject/helmit/pkg/test"
	"github.com/onosproject/rimedo-ts/test/ts"
)

func main() {
	registry.RegisterTestSuite("ts", &ts.TestSuite{})
	test.Main()
}
