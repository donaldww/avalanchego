// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package antithesis

import (
	"flag"
	"strings"

	"github.com/stretchr/testify/require"

	"github.com/ava-labs/avalanchego/tests"
	"github.com/ava-labs/avalanchego/tests/fixture/e2e"
	"github.com/ava-labs/avalanchego/tests/fixture/tmpnet"
)

// TODO(marun) Support accepting uris and chain ids as env vars

const (
	URIsKey     = "uris"
	ChainIDsKey = "chain-ids"

	FlagsName = "workload"
	EnvPrefix = "avawl"
)

type Config struct {
	URIs     []string
	ChainIDs []string
}

// Cleans up resources created by the configuration
type CleanupFunc func()

type SubnetsForNodesFunc func(nodes ...*tmpnet.Node) []*tmpnet.Subnet

func NewConfig(tc tests.TestContext, defaultNetwork *tmpnet.Network) *Config {
	return NewConfigWithSubnets(tc, defaultNetwork, nil)
}

func NewConfigWithSubnets(tc tests.TestContext, defaultNetwork *tmpnet.Network, getSubnets SubnetsForNodesFunc) *Config {
	// tmpnet configuration
	flagVars := e2e.RegisterFlags()

	var (
		uris CSV
		// Accept a list of chain IDs, assume they each belong to a separate subnet
		// TODO(marun) Revisit how chain IDs are provided when 1:n subnet:chain configuration is required.
		chainIDs CSV
	)
	flag.Var(&uris, URIsKey, "URIs of nodes that the workload can communicate with")
	flag.Var(&chainIDs, ChainIDsKey, "IDs of chains to target for testing")

	flag.Parse()

	// Use the network configuration provided
	if len(uris) != 0 {
		require.NoError(tc, awaitHealthyNodes(tc.DefaultContext(), uris), "failed to see healthy nodes")
		return &Config{
			URIs:     uris,
			ChainIDs: chainIDs,
		}
	}

	// Create a new network
	return configForNewNetwork(tc, defaultNetwork, getSubnets, flagVars)
}

// configForNewNetwork creates a new network and returns the resulting config and cleanup function.
func configForNewNetwork(tc tests.TestContext, defaultNetwork *tmpnet.Network, getSubnets SubnetsForNodesFunc, flagVars *e2e.FlagVars) *Config {
	if defaultNetwork.Nodes == nil {
		defaultNetwork.Nodes = tmpnet.NewNodesOrPanic(flagVars.NodeCount())
	}
	if defaultNetwork.Subnets == nil && getSubnets != nil {
		defaultNetwork.Subnets = getSubnets(defaultNetwork.Nodes...)
	}

	testEnv := e2e.NewTestEnvironment(tc, flagVars, defaultNetwork)

	c := &Config{}
	c.URIs = make(CSV, len(testEnv.URIs))
	for i, nodeURI := range testEnv.URIs {
		c.URIs[i] = nodeURI.URI
	}
	network := testEnv.GetNetwork()
	c.ChainIDs = make(CSV, len(network.Subnets))
	for i, subnet := range network.Subnets {
		c.ChainIDs[i] = subnet.Chains[0].ChainID.String()
	}

	return c
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
