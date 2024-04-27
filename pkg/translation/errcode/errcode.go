// Package errcode ...
package errcode

import (
	"embed"
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	instance *i18n.Bundle

	//go:embed *.toml
	fs embed.FS
)

// GetInstance returns the i18n.Bundle error codes instance.
func GetInstance() *i18n.Bundle {
	if instance == nil {
		instance = i18n.NewBundle(language.English)
		instance.RegisterUnmarshalFunc("toml", toml.Unmarshal)
		dir, err := fs.ReadDir(".")
		if err != nil {
			log.Fatalf("unable to read dir: %v", err)
		}

		for _, entry := range dir {
			info, err := entry.Info()
			if err != nil {
				log.Fatalf("unable to get info: %v", err)
			}
			b, err := fs.ReadFile(info.Name())
			if err != nil {
				log.Fatalf("unable to read file: %v", err)
			}

			instance.MustParseMessageFileBytes(b, info.Name())
		}
	}

	return instance
}

// RetrieveTranslate retrieves a message for a given key and language.
func RetrieveTranslate(key string, lang string) string {
	localizer := i18n.NewLocalizer(GetInstance(), lang, language.English.String(), language.Chinese.String())
	msgVal, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: key,
	})
	if err != nil {
		return fmt.Sprintf("$$%s", key)
	}

	return msgVal
}
