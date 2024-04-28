package language

import (
	"errors"
	"flag"
	"log/slog"

	flagUtils "github.com/madkins23/go-utils/flag"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	errNoLanguage   = errors.New("no language printer")
	languageFlags   flagUtils.StringArray
	languagePrinter *message.Printer
)

func init() {
	flag.Var(&languageFlags, "language",
		"One or more language tags to be tried, defaults to US English.")
}

// Printer returns a pre-configured message.Printer.
// The printer can execute format statements that fix numbers for the chosen language.
func Printer() *message.Printer {
	return languagePrinter
}

// Setup a message.Printer for the user or default specified language.
func Setup() error {
	languageFlags = append(languageFlags, "en_US", "en")
	for _, choice := range languageFlags {
		if tag, err := language.Parse(choice); err != nil {
			slog.Warn("Language parse failure", "choice", choice, "err", err)
		} else if languagePrinter = message.NewPrinter(tag); languagePrinter != nil {
			return nil
		}
	}
	return errNoLanguage
}
