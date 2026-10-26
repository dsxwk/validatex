package validatex

// NewRequestTranslator 创建请求翻译器
func NewRequestTranslator(request any) *Translator {
	translator := NewTranslator()
	applyRequestTranslations(translator, request)
	return translator
}

// ValidatablePointer 约束可校验的请求指针类型。
// T 为指针指向的请求类型；使用该约束的请求参数保持具体指针类型，因此可以直接与 nil 比较。
type ValidatablePointer[T any] interface {
	Validatable
	~*T
}

// ValidateRequest 使用默认翻译器校验请求,无翻译配置时翻译器为nil
func ValidateRequest[T any, P ValidatablePointer[T]](request P, scene string) Errors {
	return ValidateRequestWithTranslator(request, scene, nil)
}

// ValidateRequestWithTranslator 使用指定翻译器校验请求
func ValidateRequestWithTranslator[T any, P ValidatablePointer[T]](request P, scene string, translator *Translator) Errors {
	if request == nil {
		return nil
	}
	if translator == nil {
		translator = requestTranslator(request)
	}
	if validator, ok := any(request).(SceneValidatable); ok {
		return validator.ValidateWithScene(scene, translator)
	}
	return request.Validate(translator)
}

// requestTranslator 创建只读请求翻译器
func requestTranslator(request any) *Translator {
	_, hasTranslates := request.(TranslatesProvider)
	_, hasMessages := request.(MessagesProvider)
	if !hasTranslates && !hasMessages {
		return nil
	}
	translator := &Translator{}
	applyRequestTranslations(translator, request)
	return translator
}

// applyRequestTranslations 应用请求翻译配置
func applyRequestTranslations(translator *Translator, request any) bool {
	applied := false
	if provider, ok := request.(TranslatesProvider); ok {
		translator.Translates = provider.Translates()
		applied = true
	}
	if provider, ok := request.(MessagesProvider); ok {
		translator.Messages = provider.Messages()
		applied = true
	}
	return applied
}
