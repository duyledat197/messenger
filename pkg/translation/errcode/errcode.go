// Package errcode ...
package errcode

import (
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	instance *i18n.Bundle
)

// GetInstance returns the i18n.Bundle error codes instance.
func GetInstance() *i18n.Bundle {
	if instance == nil {
		instance = i18n.NewBundle(language.English)
		instance.RegisterUnmarshalFunc("toml", toml.Unmarshal)
		rootDir := "./"
		if _, err := os.Stat(path.Join(rootDir, "translation")); err != nil {
			if os.IsNotExist(err) {
				rootDir = "/opt/messenger"
			} else {
				slog.Error("unable to load translation file", err)
				return nil
			}
		}
		transDir := path.Join(rootDir, "translation", "errcode")

		instance.MustLoadMessageFile(path.Join(transDir, "vi.toml"))
		instance.MustLoadMessageFile(path.Join(transDir, "en.toml"))
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
