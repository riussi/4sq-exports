// Copyright © 2019 Juha Ristolainen <juha.ristolainen@iki.fi>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// These are injected at link time
var (
	version  string
	commit   string
	compiled string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "shows the application version",
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func printVersion() {
	fmt.Printf("4sq-exports %s\n\n", version)
	fmt.Printf("https://github.com/riussi/4sq-exports @%s - compiled %s\n", commit, compiled)
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
