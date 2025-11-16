package util

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	validator2 "github.com/ryuyb/fusion/internal/infrastructure/provider/validator"
	"github.com/ryuyb/fusion/internal/pkg/errors"
	"github.com/samber/lo"
)

func ParseRequestJson[T any](ctx fiber.Ctx, req *T) error {
	if err := ctx.Bind().JSON(req); err != nil {
		if errs, ok := lo.ErrorsAs[validator.ValidationErrors](err); ok {
			validationErrors := validator2.VALIDATOR.TranslateErrorsAuto(errs, ctx.Get(fiber.HeaderAcceptLanguage))
			return errors.CustomValidationError(validationErrors)
		}
		return errors.BadRequest("failed to parse request body").Wrap(err)
	}
	return nil
}
