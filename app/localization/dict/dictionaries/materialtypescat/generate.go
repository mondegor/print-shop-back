package materialtypescat

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

//go:generate gotext -srclang=en-US update -out=../../../../internal/localization/dict/dictionaries/materialtypescat/catalog.go -lang=en-US,ru-RU print-shop-back/localization/dict/dictionaries/materialtypescat
//go:generate gotext-catalog-fix -src=../../../../internal/localization/dict/dictionaries/materialtypescat/catalog.go -out=../../../../internal/localization/dict/dictionaries/materialtypescat/catalog.go

// Здесь приведены фразы используемы для локализации.
//
//nolint:unused
func list() {
	p := message.NewPrinter(language.MustParse("en-US"))

	p.Sprintf("")
}
