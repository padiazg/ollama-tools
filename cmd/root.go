/*
Copyright © 2025 Pato Diaz pato@patodiaz.io

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ollama-tools",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ollama-tools.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".ollama-tools" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".ollama-tools")
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("ot")

	setDefaults()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

}

func setDefaults() {
	// viper.SetDefault("webserver.port", 3001)
	// viper.SetDefault("webserver.adminport", 3001)
	// viper.SetDefault("webserver.tls_enabled", false)
	// viper.SetDefault("webserver.static.path", "./static")
	// viper.SetDefault("webserver.static.route", "/static")
	// viper.SetDefault("log.level", "info")
	// viper.SetDefault("log.format", "prod")
	// viper.SetDefault("internals.syncrefreshinterval", 10)
	// viper.SetDefault("internals.portalenabled", true)
	// viper.SetDefault("internals.adminenabled", true)
}

// bindEnvs creates the environment variable bindings for the given struct, also aliases for proper
// binding of environment variables and values from .env files and other structured config files.
func bindEnvs(i interface{}, parts ...string) {
	ifv := reflect.ValueOf(i)
	ift := reflect.TypeOf(i)

	// received a pointer, dereference it
	if ifv.Kind() == reflect.Ptr {
		ifv = ifv.Elem()
		ift = ift.Elem()
	}

	for x := 0; x < ift.NumField(); x++ {
		// for _, f := range reflect.VisibleFields(ift) {
		t := ift.Field(x)
		v := ifv.Field(x)

		if !t.IsExported() {
			// fmt.Printf("  not exported: %s\n", t.Name)
			continue
		}

		// fmt.Printf("field: %s, value: %s\n", t.Name, v.String())
		switch v.Kind() {
		case reflect.Struct:
			// fmt.Printf("  struct: %s\n", t.Name)
			bindEnvs(v.Interface(), append(parts, t.Name)...)

		case reflect.Ptr:
			if v.IsNil() {
				// fmt.Printf("  nil pointer: %s\n", t.Name)
				continue
			}
			// fmt.Printf("  pointer: %s\n", t.Name)
			bindEnvs(v.Interface(), append(parts, t.Name)...)

		default:
			var (
				envKey   = strings.ToUpper(strings.Join(append(parts, t.Name), "_"))
				key      = strings.Join(append(parts, t.Name), ".")
				envAlias = strings.ToLower(envKey)
			)

			// set the env binding
			if err := viper.BindEnv(key, envKey); err != nil {
				log.Fatalf("config: unable to bind env: %s", err.Error())
			}

			viper.RegisterAlias(envAlias, key)

			// fmt.Printf("  key: %s => %s => %s\n", key, envKey, envAlias)
		}
	}
}
