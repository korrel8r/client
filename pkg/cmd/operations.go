// Copyright: This file is part of korrel8r, released under https://github.com/korrel8r/korrel8r/blob/main/LICENSE

package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/korrel8r/client/pkg/swagger/client/operations"
	"github.com/korrel8r/client/pkg/swagger/models"
	"github.com/spf13/cobra"
	"k8s.io/utils/ptr"
)

// Common flags for neighbours and goals
var (
	class   string
	queries []string
	objects []string
	rules   bool

	limit                 int64
	since, until, timeout time.Duration
)

func commonFlags(cmd *cobra.Command) {
	// Start point
	cmd.Flags().StringArrayVarP(&queries, "query", "q", nil, "Query string for start objects, can be multiple.")
	cmd.Flags().StringVarP(&class, "class", "c", "", "Class for serialized start objects")
	cmd.Flags().StringArrayVarP(&objects, "object", "O", nil, "Serialized start object, can be multiple.")
	// Constraint
	cmd.Flags().Int64Var(&limit, "limit", 0, "Limit total number of results.")
	cmd.Flags().DurationVar(&timeout, "timeout", 0, "Timeout for store requests.")
	cmd.Flags().DurationVar(&since, "since", 0, "Only get results since this long ago.")
	cmd.Flags().DurationVar(&until, "until", 0, "Only get results until this long ago.")

	// Optional rules
	cmd.Flags().BoolVar(&rules, "rules", false, "Include rule information in returned graph.")
}

var domainsCmd = &cobra.Command{
	Use:   "domains",
	Short: "Get a list of domains and store configuration",
	Run: func(cmd *cobra.Command, args []string) {
		c := newClient()
		ok, err := c.Operations.GetDomains(&operations.GetDomainsParams{})
		check(err)
		NewPrinter(output.String(), os.Stdout)(ok.Payload)
	},
}

func init() {
	rootCmd.AddCommand(domainsCmd)
}

var objectsCmd = &cobra.Command{
	Use:   "objects QUERY",
	Short: "Return the list of objects for a query.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := newClient()
		ok, err := c.Operations.GetObjects(&operations.GetObjectsParams{Query: args[0]})
		check(err)
		NewPrinter(output.String(), os.Stdout)(ok.Payload)
	},
}

func init() {
	rootCmd.AddCommand(objectsCmd)
}

// Returns nil pointer for false
func ptrBool(v bool) *bool {
	if v {
		return ptr.To(v)
	}
	return nil
}

var neighboursCmd = &cobra.Command{
	Use:   "neighbours",
	Short: "Get graph of nearest neighbours",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		c := newClient()
		ok, partial, err := c.Operations.PostGraphsNeighbours(&operations.PostGraphsNeighboursParams{
			Request: &models.Neighbours{
				Depth: depth,
				Start: start(),
			},
			Rules: ptrBool(rules),
		})
		check(err)
		var payload *models.Graph
		switch {
		case ok != nil:
			payload = ok.Payload
		case partial != nil:
			fmt.Fprintln(os.Stderr, "WARNING: partial result, search timed out")
			payload = partial.Payload
		}
		NewPrinter(output.String(), os.Stdout)(payload)
	},
}

var depth int64

func init() {
	rootCmd.AddCommand(neighboursCmd)
	commonFlags(neighboursCmd)
	neighboursCmd.Flags().Int64VarP(&depth, "depth", "d", 2, "Depth of neighbourhood search.")
}

var goalsCmd = &cobra.Command{
	Use:   "goals CLASS...",
	Short: "Get graph of goal classes reachable from start",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := newClient()
		ok, partial, err := c.Operations.PostGraphsGoals(&operations.PostGraphsGoalsParams{
			Request: &models.Goals{
				Goals: args,
				Start: start(),
			},
			Rules: ptrBool(rules),
		})
		check(err)
		var payload *models.Graph
		switch {
		case ok != nil:
			payload = ok.Payload
		case partial != nil:
			fmt.Fprintln(os.Stderr, "WARNING: partial result, search timed out")
			payload = partial.Payload
		}
		NewPrinter(output.String(), os.Stdout)(payload)
	},
}

func init() {
	rootCmd.AddCommand(goalsCmd)
	commonFlags(goalsCmd)
}

var (
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Change configuration settings on the server",
		Run: func(cmd *cobra.Command, args []string) {
			config := &operations.PutConfigParams{}
			if cmd.Flags().Changed("set-verbose") {
				config.Verbose = &configVerbose
			}
			c := newClient()
			_, err := c.Operations.PutConfig(config)
			check(err)
		},
	}

	configVerbose int64
)

func init() {
	configCmd.Flags().Int64Var(&configVerbose, "set-verbose", 0, "Set verbose level for logging")
	rootCmd.AddCommand(configCmd)
}

func start() *models.Start {
	return &models.Start{
		Class:      class,
		Constraint: constraint(),
		Objects:    objects,
		Queries:    queries,
	}
}

func constraint() *models.Constraint {
	c := &models.Constraint{Limit: limit}
	if timeout > 0 {
		c.Timeout = timeout.String()
	}
	now := time.Now()
	if since > 0 {
		c.Start = ptr.To(strfmt.DateTime(now.Add(-since)))
	}
	if until > 0 {
		c.End = ptr.To(strfmt.DateTime(now.Add(-until)))
	}
	return c
}
