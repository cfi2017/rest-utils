package util

import (
	"fmt"
	"github.com/spf13/pflag"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// InitialiseConfig initialises a generic config.
//
// Any registered flags are available as environment variables under the `[APP]` prefix:
// e.g. db.username => BACKEND_DB.USERNAME=backend
//
// Any registered flags are also available as command line arguments:
// e.g. db.username => ./app --db.username=backend
//
// Any registered flags are also available within toml, yaml, json or xml configuration files:
// e.g. db.username => backend.toml
// [db]
// username = app
//
// The config library looks for config files in the following paths:
// - /etc/[app]/
// - $HOME/.[app]/
// - ./config/
// - . (working directory)
//
// The name of the config file ([app].toml) depends on the argument passed to InitialiseConfig:
// e.g. "backend" => /etc/faceit/backend.toml, backend.yaml, backend.json...
func InitialiseConfig(name string) {

	// look for env variables in the format "[APP]_PORT=1338"
	viper.SetEnvPrefix(strings.ToUpper(name))

	// look for config files with name name.yml, name.toml, name.json...
	viper.SetConfigName(name)

	// ... in these folders
	viper.AddConfigPath(fmt.Sprintf("/etc/%s", name))
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s", name))
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".") // working directory

	// parse flags from process arg list
	pflag.Parse()

	// bind parsed flags to config library
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic(err)
	}

	// check for environment variables now
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// try to find and read config file now
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
			panic(err)
		}
	}

	// watch config file for changes
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		GetS().Info("Config file changed:", e.Name)
	})
	if used := viper.ConfigFileUsed(); used != "" {
		GetS().Info("Using config file:", used)
	}

}
