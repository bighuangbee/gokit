package i18n

import (
	"context"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Options struct {
	// DefaultLanguage *language.Tag
	Paths []string
}

// New
// paths 文件目录地址，文件命名规则: {name}.zh.toml, {name}.en.toml, zh和en固定
func New(opt Options) *i18n.Bundle {
	bundle := i18n.NewBundle(language.Chinese)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	for _, path := range opt.Paths {
		bundle.MustLoadMessageFile(path)
	}
	return bundle
}

type ILocalizer interface {
	Tr(string, ...interface{}) (string, error)
}

type localizer struct {
	l *i18n.Localizer
}

func newLocalizer(bundle *i18n.Bundle, lang string) localizer {
	return localizer{
		l: i18n.NewLocalizer(bundle, lang),
	}
}

func (t localizer) Tr(msgID string, args ...interface{}) (string, error) {
	size := len(args)
	var data map[string]interface{}
	if size > 0 {
		data = make(map[string]interface{}, size)
		for i := 1; i <= size; i++ {
			data["p"+strconv.Itoa(i)] = args[i-1]
		}
	}
	return t.l.Localize(&i18n.LocalizeConfig{
		MessageID:    msgID,
		TemplateData: data,
	})
}

func NewContext(ctx context.Context, value ILocalizer) context.Context {
	return context.WithValue(ctx, localizer{}, value)
}

func FromContext(ctx context.Context) (ILocalizer, bool) {
	l, ok := ctx.Value(localizer{}).(ILocalizer)
	return l, ok
}

// Tr 翻译
func Tr(ctx context.Context, msgID string, args ...interface{}) string {
	l, ok := FromContext(ctx)
	msg := msgID
	if ok {
		if result, err := l.Tr(msgID, args...); err == nil {
			msg = result
		}
	}
	return msg
}
