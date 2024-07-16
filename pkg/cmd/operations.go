// Copyright: This file is part of korrel8r, released under https://github.com/korrel8r/korrel8r/blob/main/LICENSE

package cmd

import (
	"os"

	"github.com/korrel8r/client/pkg/swagger/client/operations"
	"github.com/korrel8r/client/pkg/swagger/models"
	"github.com/spf13/cobra"
)

// Common flags
var (
	queries []string
	class   string
	objects []string
	depth   int64
	rules   bool
)

func makeStart() *models.Start {
	return &models.Start{
		Class:      class,
		Constraint: nil, // TODO support for constraints.
		Objects:    objects,
		Queries:    queries,
	}
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

func init() { rootCmd.AddCommand(domainsCmd) }

var (
	objectsCmd = &cobra.Command{
		Use:   "objects [QUERY]",
		Short: "Return the list of objects for a query.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c := newClient()
			ok, err := c.Operations.GetObjects(&operations.GetObjectsParams{Query: args[0]})
			check(err)
			NewPrinter(output.String(), os.Stdout)(ok.Payload)
		},
	}
)

func init() {
	rootCmd.AddCommand(objectsCmd)
}

var neighboursCmd = &cobra.Command{
	Use:   "neighbours [FLAGS]",
	Short: "Get graph of nearest neighbours",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		c := newClient()
		ok, err := c.Operations.PostGraphsNeighbours(&operations.PostGraphsNeighboursParams{
			Request: &models.Neighbours{
				Depth: depth,
				Start: makeStart(),
			},
			Rules: &rules,
		})
		check(err)
		NewPrinter(output.String(), os.Stdout)(ok.Payload)
	},
}

var goalsCmd = &cobra.Command{
	Use:   "goals [FLAGS] CLASS...",
	Short: "Get graph of goal classes reachable from start",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := newClient()
		ok, err := c.Operations.PostGraphsGoals(&operations.PostGraphsGoalsParams{
			Request: &models.Goals{
				Goals: args,
				Start: makeStart(),
			},
			Rules: &rules,
		})
		check(err)
		NewPrinter(output.String(), os.Stdout)(ok.Payload)
	},
}

func commonFlags(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
	cmd.Flags().StringArrayVarP(&queries, "query", "q", nil, "Query string for start objects, can be multiple.")
	cmd.Flags().StringVarP(&class, "class", "c", "", "Class for serialized start objects")
	cmd.Flags().StringArrayVarP(&objects, "object", "0", nil, "Serialized start object, can be multiple.")
	cmd.Flags().BoolVar(&rules, "rules", false, "Include per-rule information in returned graph.")
}

func init() {
	commonFlags(neighboursCmd)
	neighboursCmd.Flags().Int64Var(&depth, "depth", 2, "Depth of neighbourhood search.")
	commonFlags(goalsCmd)
}
