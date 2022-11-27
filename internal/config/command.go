package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ChrisWiegman/kana-cli/pkg/console"

	"github.com/aquasecurity/table"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
)

func (c *Config) GetGlobalSetting(md *cobra.Command, args []string) (string, error) {

	if !c.Global.viper.IsSet(args[0]) {
		return "", fmt.Errorf("invalid setting. Please enter a valid key to get")
	}

	return c.Global.viper.GetString(args[0]), nil
}

func (c *Config) ListConfig() {

	t := table.New(os.Stdout)

	t.SetHeaders("Setting", "Global Value", "Local Value")

	t.AddRow("admin.email", console.Bold(c.Global.AdminEmail))
	t.AddRow("admin.password", console.Bold(c.Global.AdminPassword))
	t.AddRow("admnin.username", console.Bold(c.Global.AdminUsername))
	t.AddRow("local", console.Bold(strconv.FormatBool(c.Global.Local)), console.Bold(strconv.FormatBool(c.Local.Local)))
	t.AddRow("php", console.Bold(c.Global.PHP), console.Bold(c.Local.PHP))
	t.AddRow("type", console.Bold(c.Global.Type), console.Bold(c.Local.Type))
	t.AddRow("xdebug", console.Bold(strconv.FormatBool(c.Global.Xdebug)), console.Bold(strconv.FormatBool(c.Local.Local)))

	boldPlugins := []string{}

	for _, plugin := range c.Local.Plugins {
		boldPlugins = append(boldPlugins, console.Bold(plugin))
	}

	plugins := console.Bold(strings.Join(boldPlugins, "\n"))

	t.AddRow("plugins", "", plugins)

	t.Render()
}

func (c *Config) SetGlobalSetting(md *cobra.Command, args []string) error {

	if !c.Global.viper.IsSet(args[0]) {
		return fmt.Errorf("invalid setting. Please enter a valid key to set")
	}

	validate := validator.New()
	var err error

	switch args[0] {
	case "local", "xdebug":
		err = validate.Var(args[1], "boolean")
		if err != nil {
			return err
		}
		boolVal, err := strconv.ParseBool(args[1])
		if err != nil {
			return err
		}
		c.Global.viper.Set(args[0], boolVal)
		return c.Global.viper.WriteConfig()
	case "php":
		if !isValidString(args[1], validPHPVersions) {
			err = fmt.Errorf("please choose a valid php version")
		}
	case "type":
		if !isValidString(args[1], validTypes) {
			err = fmt.Errorf("please choose a valid project type")
		}
	case "admin.email":
		err = validate.Var(args[1], "email")
	case "admin.password":
		err = validate.Var(args[1], "alphanumunicode")
	case "admin.username":
		err = validate.Var(args[1], "alpha")
	default:
		err = validate.Var(args[1], "boolean")
	}

	if err != nil {
		return err
	}

	c.Global.viper.Set(args[0], args[1])

	return c.Global.viper.WriteConfig()
}
