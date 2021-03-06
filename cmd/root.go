package cmd

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/qlik-oss/corectl/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var headersMap = make(map[string]string)
var explicitConfigFile = ""
var version = ""
var headers http.Header
var rootCtx = context.Background()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Hidden:                 true,
	Use:                    "corectl",
	Short:                  "",
	Long:                   `corectl contains various commands to interact with the Qlik Associative Engine. See respective command for more information`,
	DisableAutoGenTag:      true,
	BashCompletionFunction: bashCompletionFunc,

	Annotations: map[string]string{
		"x-qlik-stability": "experimental",
	},

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		// if help, version or generate-docs command, no prerun is needed.
		if strings.Contains(ccmd.Use, "help") || ccmd.Use == "generate-docs" || ccmd.Use == "generate-spec" || ccmd.Use == "version" {
			return
		}
		internal.ReadConfigFile(explicitConfigFile)

		if len(headersMap) == 0 {
			headersMap = viper.GetStringMapString("headers")
		}
		headers = make(http.Header, 1)
		for key, value := range headersMap {
			headers.Set(key, value)
		}
	},

	Run: func(ccmd *cobra.Command, args []string) {
		ccmd.HelpFunc()(ccmd, args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(mainVersion string) {
	version = mainVersion
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func patchRootCommandUsageTemplate() {
	var originalUsageSnippet = `Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}`

	var rootSnippetMainSection = `App Building Commands:{{range .Commands}}{{if (and (or .IsAvailableCommand (eq .Name "help")) (eq (index .Annotations "command_category") "build"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}

App Analysis Commands:{{range .Commands}}{{if (and .IsAvailableCommand (eq (index .Annotations "command_category") ""))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}

Advanced Commands:{{range .Commands}}{{if (and (or .IsAvailableCommand (eq .Name "help")) (eq (index .Annotations "command_category") "sub"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}

Other Commands:{{range .Commands}}{{if (and (or .IsAvailableCommand (eq .Name "help")) (or (eq (index .Annotations "command_category") "other") (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}`

	var newUsageSnippet = `{{if (eq .Name "corectl")}}` + rootSnippetMainSection + `{{else}}` + originalUsageSnippet + "{{end}}"

	var patchedUsageTemplate = strings.Replace(rootCmd.UsageTemplate(), originalUsageSnippet, newUsageSnippet, 1)
	rootCmd.SetUsageTemplate(patchedUsageTemplate)
}

func init() {

	// Common commands
	rootCmd.AddCommand(getTablesCmd)
	rootCmd.AddCommand(getFieldsCmd)
	rootCmd.AddCommand(getAssociationsCmd)
	rootCmd.AddCommand(getKeysCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(evalCmd)
	rootCmd.AddCommand(reloadCmd)
	rootCmd.AddCommand(getValuesCmd)
	rootCmd.AddCommand(getMetaCmd)

	// Subcommands
	rootCmd.AddCommand(measureCmd)
	rootCmd.AddCommand(dimensionCmd)
	rootCmd.AddCommand(objectCmd)
	rootCmd.AddCommand(connectionCmd)
	rootCmd.AddCommand(scriptCmd)
	rootCmd.AddCommand(appCmd)

	// Other
	rootCmd.AddCommand(catwalkCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(generateDocsCmd)
	rootCmd.AddCommand(generateSpecCmd)

	initGlobalFlags(rootCmd.PersistentFlags())
	patchRootCommandUsageTemplate()

}
