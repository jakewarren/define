package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	dictionary "github.com/jakewarren/define"
	"github.com/mitchellh/go-homedir"
	"github.com/ogier/pflag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func main() {

	config := pflag.StringP("config", "c", "~/.define.conf.json", "path to config file")
	disableColor := pflag.Bool("disable-color", false, "disable color output")
	pflag.Parse()

	if *disableColor || envDisableColor() {
		color.NoColor = true
	}

	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	if pflag.NArg() < 1 {
		log.Error().Msg("No word provided")
		pflag.Usage()
		os.Exit(1)
	}

	var (
		client        = http.Client{Timeout: 19 * time.Second}
		appID, appKey string
	)

	// check if the key is present in ENV variables. if not read in config file
	if _, ok := os.LookupEnv("OXFORD_DICTIONARY_APP_ID"); ok {
		appID = os.Getenv("OXFORD_DICTIONARY_APP_ID")
		appKey = os.Getenv("OXFORD_DICTIONARY_APP_KEY")
	} else {
		viper.SetConfigFile(tryExpandPath(*config))

		err := viper.ReadInConfig()
		if err != nil {
			log.Error().Err(err).Msg("error reading config")
		}

		appID = viper.GetString("OxfordDictionary.AppID")
		appKey = viper.GetString("OxfordDictionary.AppKey")
	}

	c := dictionary.New(client, appID, appKey)

	r, err := c.Define(pflag.Arg(0))

	if err != nil {
		log.Fatal().Err(err).Msg("error retrieving definition")
	}

	prettyPrint(r)
}

func envDisableColor() bool {
	// check for the existence of NO_COLOR to satisfy the nocolor standard http://no-color.org
	_, exists := os.LookupEnv("NO_COLOR")
	if exists {
		return true
	}
	// disable color if terminal is set to dumb
	val, _ := os.LookupEnv("TERM")
	return val == "DUMB"
}

func prettyPrint(resp *dictionary.Results) {

	for _, r := range resp.Results {
		for _, lexEntry := range r.LexicalEntries {

			// print part of speech
			fmt.Println(color.YellowString("%s\n", lexEntry.RenderLexicalCategory()))
			for _, e := range lexEntry.Entries {
				for i, sense := range e.Senses {
					headerFmt := color.New(color.FgHiMagenta, color.Bold)
					exampleFmt := color.New(color.Italic)
					headerFmt.Println(fmt.Sprintf("%d.", i+1), rejoin(sense.Definitions[0], "\n", ", "))

					if sense.HasExamples() {
						exampleFmt.Println("\t", "'"+rejoin(sense.RenderExamples(), "\n", ", ")+"'")
					}

					for j, subsense := range sense.Subsenses {
						subsenseFmt := color.New(color.FgHiCyan, color.Bold)
						subsenseFmt.Printf("\n   %s - %s%s\n", fmt.Sprintf("%d.%d", i+1, j+1), subsense.RenderTags(), subsenseFmt.Sprint(subsense.Definitions[0]))
						if subsense.HasExamples() {
							exampleFmt.Println("\t", "'"+rejoin(subsense.RenderExamples(), "\n", ", ")+"'")
						}

					}

					fmt.Println()
				}

				// print the origin
				if e.HasEtymology() {
					fmt.Println(color.YellowString("%s", "ORIGIN"))
					fmt.Println("\t", e.RenderEtymology())
				}

			}

		}

	}
}

// Helper function to split a string on an old delimeter and re-join the string
// slice with a new one.
func rejoin(input, old, new string) string {
	input = strings.TrimRight(input, old)
	tmp := strings.Split(input, old)
	return strings.Join(tmp, new)
}

// tryExpandPath attempts to expand a given path and returns the expanded path
// if successful. Otherwise, if expansion failed, the original path is returned.
func tryExpandPath(path string) string {
	if expanded, err := homedir.Expand(path); nil == err {
		path = expanded
	}

	return path
}
