/*
Copyright © 2023 Kun <ev.sin@hotmal.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"gopw-crawler/cmd/app"
	"log"

	"github.com/spf13/cobra"
)

func NewDailyPowerQueryCmd() *cobra.Command {
	opt := app.NewOptions()

	dailyPowerQueryCmd := &cobra.Command{
		Use:   "dailyPowerQuery",
		Short: "Batch download daily power query sheets",
		Long:  `Batch download daily power query sheets`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := opt.Complete(); err != nil {
				log.Fatalf("Failed to complete options: %v", err)
			}
			if err := opt.Run(); err != nil {
				log.Fatalf("Failed to run: %v", err)
			}
		},
	}
	if err := opt.Init(); err != nil {
		log.Fatalf("Failed to init options: %v", err)
	}
	opt.AddAllFlags(dailyPowerQueryCmd.Flags())
	return dailyPowerQueryCmd
}

func init() {

	rootCmd.AddCommand(NewDailyPowerQueryCmd())

}
