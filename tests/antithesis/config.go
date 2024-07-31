// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package antithesis

import (
	"errors"
	"flag"
	"strings"
	"time"

	"github.com/ava-labs/avalanchego/tests"
	"github.com/ava-labs/avalanchego/tests/fixture/e2e"
	"github.com/ava-labs/avalanchego/tests/fixture/tmpnet"
)

const (
	URIsKey     = "uris"
	ChainIDsKey = "chain-ids"

	FlagsName = "workload"
	EnvPrefix = "avawl"
)

var (
	errNoURIs      = errors.New("at least one URI must be provided")
	errNoArguments = errors.New("no arguments")
)

type Config struct {
	URIs          []string
	ChainIDs      []string
	ReuseNetwork  bool
	ShutdownDelay time.Duration
}

// TODO(marun) Revisit whether errors should be propagated instead of failing directly
func NewConfig(defaultNetwork *tmpnet.Network, getSubnets func(nodes ...*tmpnet.Node) []*tmpnet.Subnet) *Config {
	// tmpnet configuration
	flagVars := e2e.RegisterFlags()

	var (
		uris     CSV
		chainIDs CSV
	)
	// TODO(marun) Support the case of accepting a local api uri e.g. primary.LocalAPIURI
	flag.Var(&uris, URIsKey, "URIs of nodes that the workload can communicate with")
	flag.Var(&chainIDs, ChainIDsKey, "IDs of chains to target for testing")

	flag.Parse()

	if len(uris) == 0 {
		// No URIs provided, default to tmpnet
		defaultNetwork.Nodes = tmpnet.NewNodesOrPanic(flagVars.NodeCount())
		defaultNetwork.Subnets = getSubnets(defaultNetwork.Nodes...)
		tc := tests.NewTestContext()
		testEnv := e2e.NewTestEnvironment(tc, flagVars, defaultNetwork)
		uris = make(CSV, len(testEnv.URIs))
		for i, nodeURI := range testEnv.URIs {
			uris[i] = nodeURI.URI
		}
		network := testEnv.GetNetwork()
		chainIDs = make(CSV, len(network.Subnets))
		for i, subnet := range network.Subnets {
			chainIDs[i] = subnet.Chains[0].ChainID.String()
		}
	}

	return &Config{
		URIs:          uris,
		ChainIDs:      chainIDs,
		ReuseNetwork:  flagVars.ReuseNetwork(),
		ShutdownDelay: flagVars.NetworkShutdownDelay(),
	}
}

// CSV is a custom type that implements the flag.Value interface
type CSV []string

// String returns the string representation of the CSV type
func (c *CSV) String() string {
	return strings.Join(*c, ",")

}

// Set splits the input string by commas and sets the CSV type
func (c *CSV) Set(value string) error {
	*c = strings.Split(value, ",")
	return nil
}
