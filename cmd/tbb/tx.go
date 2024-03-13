package main

import (
	"github.com/spf13/cobra"
)

const flagFrom = "from"
const flagTo = "to"
const flagValue = "value"
const flagData = "data"

func txCmd() *cobra.Command {
	var txsCmd = &cobra.Command{
		Use:   "tx",
		Short: "Interact with txs (add...).",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return incorrectUsageErr()
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	txsCmd.AddCommand(txAddCmd())

	return txsCmd
}

func txAddCmd() *cobra.Command {
	var txAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Add a tx.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return incorrectUsageErr()
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	txAddCmd.Flags().String(flagFrom, "", "From what account to send the tx")
	txAddCmd.Flags().String(flagTo, "", "To what account to send the tx")
	txAddCmd.Flags().Int(flagValue, 0, "Value to send with the tx")
	txAddCmd.Flags().String(flagData, "", "Data to send with the tx")

	return txAddCmd
}
