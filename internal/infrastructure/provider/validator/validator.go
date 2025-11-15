package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag,omitempty"`
}

type Validator struct {
	uni      *ut.UniversalTranslator
	validate *validator.Validate
}

func NewValidator(logger *zap.Logger) *Validator {
	enLocales := en.New()
	universalTranslator := ut.New(enLocales, enLocales, zh.New())

	validate := validator.New()
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	registerTranslations(validate, universalTranslator, logger)

	return &Validator{
		uni:      universalTranslator,
		validate: validate,
	}
}

func registerTranslations(validate *validator.Validate, universalTranslator *ut.UniversalTranslator, logger *zap.Logger) {
	enTrans, _ := universalTranslator.GetTranslator("en")
	if err := enTranslations.RegisterDefaultTranslations(validate, enTrans); err != nil {
		logger.Fatal("failed to register english translations", zap.Error(err))
		panic(err)
	}
	zhTrans, _ := universalTranslator.GetTranslator("zh")
	if err := zhTranslations.RegisterDefaultTranslations(validate, zhTrans); err != nil {
		logger.Fatal("failed to register zh translations", zap.Error(err))
		panic(err)
	}
}

func (v *Validator) Validate(i any) error {
	return v.validate.Struct(i)
}

func (v *Validator) ValidateVar(field any, tag string) error {
	return v.validate.Var(field, tag)
}

func (v *Validator) TranslateErrors(err error, locale string) []ValidationError {
	trans, _ := v.uni.GetTranslator(locale)

	var validationErrors []ValidationError

	if errs, ok := lo.ErrorsAs[validator.ValidationErrors](err); ok {
		for _, e := range errs {
			validationErrors = append(validationErrors, ValidationError{
				Field:   e.Field(),
				Message: e.Translate(trans),
				Tag:     e.Tag(),
			})
		}
	}

	return validationErrors
}

func (v *Validator) TranslateErrorsAuto(err error, acceptLanguage string) []ValidationError {
	locale := parseLocale(acceptLanguage)
	return v.TranslateErrors(err, locale)
}

func parseLocale(locale string) string {
	locale = strings.ToLower(locale)
	switch {
	case strings.HasPrefix(locale, "zh"):
		return "zh"
	case strings.HasPrefix(locale, "en"):
		return "en"
	default:
		return "en"
	}
}
