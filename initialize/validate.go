package initialize

import (
	"fmt"
	"go-gin-rest-api/pkg/global"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

// 初始化校验器
func Validate(locale string) {
	// 修改gin框架中的Validator引擎属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册一个获取json tag的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器
		// 第一个参数是备用（fallback）的语言环境
		// 后面的参数是应该支持的语言环境（支持多个）
		// uni := ut.New(zhT, zhT) 也是可以的
		uni := ut.New(enT, zhT, enT)
		// locale 通常取决于 http 请求头的 'Accept-Language'
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		var ok bool
		global.Translator, ok = uni.GetTranslator(locale)
		if !ok {
			global.Log.Error(fmt.Sprintf("初始化validator.v10校验器 uni.GetTranslator(%s) 失败", locale))
		}
		var err error
		// 注册翻译器
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, global.Translator)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, global.Translator)
		default:
			err = zhTranslations.RegisterDefaultTranslations(v, global.Translator)
		}
		if err != nil {
			global.Log.Error("初始化validator.v10校验器失败", err)
		}
	}
	global.Log.Info("初始化validator.v10校验器完成")
}
