package main

import (
	"github.com/RIMEDO-Labs/rimedo-ts/test/ts"
	"github.com/onosproject/helmit/pkg/registry"
	"github.com/onosproject/helmit/pkg/test"
)

func main() {
	registry.RegisterTestSuite("ts", &ts.TestSuite{})
	test.Main()
}


