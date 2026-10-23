package validatex

// NewRequestTranslator 创建请求翻译器
func NewRequestTranslator(request any) *Translator {
	translator := NewTranslator()
	applyRequestTranslations(translator, request)
	return translator
}

// ValidateRequest 使用默认翻译器校验请求,无翻译配置时翻译器为nil
func ValidateRequest[T Validatable](request T, scene string) Errors {
	return ValidateRequestWithTranslator(request, scene, nil)
}

// ValidateRequestWithTranslator 使用指定翻译器校验请求
func ValidateRequestWithTranslator[T Validatable](request T, scene string, translator *Translator) Errors {
	if any(request) == nil {
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
