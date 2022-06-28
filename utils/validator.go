package utils

import (
	"errors"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

type Validator struct {
	Trans     ut.Translator
	Validator *validator.Validate
	Some      interface{}
	Err       error
}

func NewValidator(Some interface{}) *Validator {
	uni := ut.New(zh.New())
	trans, _ := uni.GetTranslator("zh")
	validate := validator.New()
	//验证器注册翻译器
	return &Validator{
		Trans:     trans,
		Validator: validate,
		Some:      Some,
		Err:       zh_translations.RegisterDefaultTranslations(validate, trans),
	}
}

func (validate *Validator) IsOk() (bool, error) {

	if validate.Err != nil {
		return true, validate.Err
	}

	err := validate.Validator.Struct(validate.Some)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return true, errors.New(err.Translate(validate.Trans))
		}
	}
	return false, nil
}
