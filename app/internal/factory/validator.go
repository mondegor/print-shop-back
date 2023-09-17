package factory

import (
    "print-shop-back/config"
    "print-shop-back/internal/controller/view"

    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrview"
)

func NewValidator(cfg *config.Config, logger mrcore.Logger) (mrcore.Validator, error) {
    logger.Info("Create and init data validator")

    validator := mrview.NewValidator()

    err := validator.Register("article", view.ValidateArticle)

    if err != nil {
        return nil, err
    }

    err = validator.Register("variable", view.ValidateVariable)

    if err != nil {
        return nil, err
    }

    return validator, nil
}
