package controller

import (
	"bluebell/models"
	"fmt"
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

// 定义一个全局翻译器
var trans ut.Translator

// InitTrans 初始化翻译器
func InitTrans(locale string) (err error) {
	// 修改gin框架中的Validator引擎属性，实现自定制
	/*这行代码从 Gin 框架中获取当前的验证器引擎，确保它是 go-playground/validator 类型。
	Gin 中的 binding.Validator 是用于请求数据绑定和校验的工具。*/
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		// 注册一个获取json tag的自定义方法
		/*注册了一个自定义的标签解析函数，使得 json 标签可以在结构体字段中正确地获取，
		默认情况下，Go 通过 json 标签来解析结构体字段。"-" 表示忽略该字段。*/
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		// 为SignUpParam注册自定义校验方法
		v.RegisterStructValidation(SignUpParamStructLevelValidation, models.SignUpParam{})

		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器

		// 第一个参数是备用（fallback）的语言环境
		// 后面的参数是应该支持的语言环境（支持多个）
		// uni := ut.New(zhT, zhT) 也是可以的
		uni := ut.New(enT, zhT, enT)

		// locale 通常取决于 http 请求头的 'Accept-Language'
		var ok bool
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}

		// 注册翻译器
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		}
		return
	}
	return
}

// 去除提示信息中的结构体名称
/*移除结构体名称前缀，从而清理校验错误信息。
比如，如果有一个字段 User.Name，那么通过此函数，会将其转换为 Name，方便返回更简洁的错误信息*/
func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}

// 自定义SignUpParam结构体校验函数
/*这是一个自定义的结构体校验函数，主要用于校验注册参数中的 Password 和 RePassword 是否一致。
具体校验规则是：如果 Password 和 RePassword 不相同，则通过 sl.ReportError 报告错误，提示用户密码和确认密码不一致。*/
func SignUpParamStructLevelValidation(sl validator.StructLevel) {
	su := sl.Current().Interface().(models.SignUpParam)

	if su.Password != su.RePassword {
		// 输出错误提示信息，最后一个参数就是传递的param
		sl.ReportError(su.RePassword, "re_password", "RePassword", "eqfield", "password")
	}
}
