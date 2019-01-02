package main_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	"os"
	"testing"
)

func TestApplication(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter(os.Getenv("CI_REPORT"))
	RunSpecsWithDefaultAndCustomReporters(t, "Gopher Test Suite", []Reporter{junitReporter})
}
