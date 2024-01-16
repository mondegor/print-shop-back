package dictionaries

import . "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrLaminateTypeNotFound = NewFactory(
		"errDictionariesLaminateTypeNotFound", ErrorKindUser, "laminate type with ID={{ .id }} not found")

	FactoryErrPaperColorNotFound = NewFactory(
		"errDictionariesPaperColorNotFound", ErrorKindUser, "paper color with ID={{ .id }} not found")

	FactoryErrPaperFactureNotFound = NewFactory(
		"errDictionariesPaperFactureNotFound", ErrorKindUser, "paper facture with ID={{ .id }} not found")

	FactoryErrPrintFormatNotFound = NewFactory(
		"errDictionariesPrintFormatNotFound", ErrorKindUser, "print format with ID={{ .id }} not found")
)
