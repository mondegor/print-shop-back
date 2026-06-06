package msgcat

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

//go:generate gotext -srclang=en-US update -out=../../../internal/localization/dict/msgcat/catalog.go -lang=en-US,ru-RU print-shop-back/localization/dict/msgcat
//go:generate gotext-catalog-fix -src=../../../internal/localization/dict/msgcat/catalog.go -out=../../../internal/localization/dict/msgcat/catalog.go

// Здесь приведены фразы используемы для локализации.
//
//nolint:unused
func list() {
	p := message.NewPrinter(language.MustParse("en-US"))

	p.Sprintf("Confirm your identity to sign in")
}
