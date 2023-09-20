package initializers

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

type Locale string

const (
	ZH_Locale = "zh"
	EN_Locale = "en"
)

var (
	once     sync.Once
	trans    ut.Translator
	validate *validator.Validate
)

// InitValidator 初始化验证器，仅会执行一次，locale 默认是 zh
func InitValidator(locales ...Locale) (ut.Translator, *validator.Validate) {
	// 仅第一次调用的时候才会调用
	once.Do(func() {
		locale := "zh"
		if len(locales) > 0 {
			locale = string(locales[0])
		}

		// 设置翻译语言支持
		en := en.New()
		zh := zh.New()
		uni := ut.New(en, en, zh)

		// 这里不需要做额外的处理，找不到zh就会使用会退的en
		var ok bool
		trans, ok = uni.GetTranslator(locale)
		fmt.Printf("->>> uni.GetTranslator(%q) ok = %t \n", locale, ok)

		validate = validator.New()

		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			if name == "" {
				name = fld.Name
			}
			return name + " "
		})

		switch locale {
		case "zh":
			zh_translations.RegisterDefaultTranslations(validate, trans)
		default:
			en_translations.RegisterDefaultTranslations(validate, trans)
		}
	})

	return trans, validate
}
