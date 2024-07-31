// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package tests

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	"github.com/onsi/ginkgo/v2/formatter"

	"github.com/stretchr/testify/require"

	"github.com/ava-labs/avalanchego/wallet/subnet/primary/common"
)

// TODO(marun) Choose a better name
type OtherTestContext struct {
	defaultTimeout time.Duration
}

func NewTestContext() *OtherTestContext {
	return &OtherTestContext{
		// TODO(marun) The default value should probably be centralized
		defaultTimeout: 2 * time.Minute,
	}
}

func (tc *OtherTestContext) Errorf(format string, args ...interface{}) {
	log.Printf("error: "+format, args...)
}

func (tc *OtherTestContext) FailNow() {
	log.Fatal("FailNow called")
}

func (tc *OtherTestContext) GetWriter() io.Writer {
	return os.Stdout
}

func (tc *OtherTestContext) DeferCleanup(args ...interface{}) {
	// TODO(marun) Register functions and ensure they are called
}

func (tc *OtherTestContext) By(text string, callback ...func()) {
	// TODO(marun)
}

// Outputs to stdout.
//
// Examples:
//
//   - Out("{{green}}{{bold}}hi there %q{{/}}", "aa")
//   - Out("{{magenta}}{{bold}}hi therea{{/}} {{cyan}}{{underline}}b{{/}}")
//
// See https://github.com/onsi/ginkgo/blob/v2.0.0/formatter/formatter.go#L52-L73
// for an exhaustive list of color options.
func (tc *OtherTestContext) Outf(format string, args ...interface{}) {
	s := formatter.F(format, args...)
	// Use GinkgoWriter to ensure that output from this function is
	// printed sequentially within other test output produced with
	// GinkgoWriter (e.g. `STEP:...`) when tests are run in
	// parallel. ginkgo collects and writes stdout separately from
	// GinkgoWriter during parallel execution and the resulting output
	// can be confusing.
	log.Print(s)
}

// Helper simplifying use of a timed context by canceling the context on ginkgo teardown.
func (tc *OtherTestContext) ContextWithTimeout(duration time.Duration) context.Context {
	return ContextWithTimeout(tc, duration)
}

// Helper simplifying use of a timed context configured with the default timeout.
func (tc *OtherTestContext) DefaultContext() context.Context {
	return DefaultContext(tc)
}

// Helper simplifying use via an option of a timed context configured with the default timeout.
func (tc *OtherTestContext) WithDefaultContext() common.Option {
	return WithDefaultContext(tc)
}

// Re-implementation of testify/require.Eventually that is compatible with ginkgo. testify's
// version calls the condition function with a goroutine and ginkgo assertions don't work
// properly in goroutines.
func (tc *OtherTestContext) Eventually(condition func() bool, waitFor time.Duration, tick time.Duration, msg string) {
	require.Eventually(tc, condition, waitFor, tick, msg)
}
