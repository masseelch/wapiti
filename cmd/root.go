/*
Copyright © 2021 MasseElch info@masseelch.de

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"github.com/masseelch/wapiti/wapiti"
	"github.com/masseelch/wapiti/wapiti/config"
	"github.com/spf13/cobra"
	"log"
)

var cfg *config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wapiti",
	Short: "Interactive cli to create ent schemas",
	Run: func(cmd *cobra.Command, args []string) {
		w, err := wapiti.New(cfg)
		fatalOnErr(err)
		fatalOnErr(w.Run())
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cfg = new(config.Config)
	rootCmd.Flags().StringVar(&cfg.SchemaPath, "schema", "ent/schema", "/path/to/schema/dir")
}

func fatalOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
