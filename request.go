package validatex

// NewRequestTranslator 创建已编译的请求翻译器快照
func NewRequestTranslator[T any, P PointerTo[T]](request P) *Translator {
	translator := NewTranslator()
	if request != nil {
		applyRequestTranslations(request, translator)
	}
	return translator.Compile()
}

// PointerTo 约束具体指针类型,T为指针指向的值类型
type PointerTo[T any] interface {
	~*T
}

// ValidatablePointer 约束可校验的请求指针类型
// T为指针指向的请求类型
type ValidatablePointer[T any] interface {
	Validatable
	PointerTo[T]
}

// SceneValidatablePointer 约束场景校验请求指针类型
// T为指针指向的请求类型
type SceneValidatablePointer[T any] interface {
	SceneValidatable
	PointerTo[T]
}

// ValidateRequest 使用默认翻译器校验请求
func ValidateRequest[T any, P ValidatablePointer[T]](request P) Errors {
	return ValidateRequestWithTranslator(request, nil)
}

// ValidateRequestWithTranslator 使用指定翻译器校验请求
func ValidateRequestWithTranslator[T any, P ValidatablePointer[T]](request P, translator *Translator) Errors {
	if request == nil {
		return nil
	}
	if translator == nil {
		translator = requestTranslator(request)
	}
	return request.Validate(translator)
}

// ValidateSceneRequest 使用默认翻译器校验场景请求
func ValidateSceneRequest[T any, P SceneValidatablePointer[T]](request P, scene string) Errors {
	return ValidateSceneRequestWithTranslator(request, scene, nil)
}

// ValidateSceneRequestWithTranslator 使用指定翻译器校验场景请求
func ValidateSceneRequestWithTranslator[T any, P SceneValidatablePointer[T]](request P, scene string, translator *Translator) Errors {
	if request == nil {
		return nil
	}
	if translator == nil {
		translator = requestTranslator(request)
	}
	return request.ValidateWithScene(scene, translator)
}

// requestTranslator 创建已编译的只读请求翻译器
func requestTranslator(request any) *Translator {
	_, providesTranslates := request.(TranslatesProvider)
	_, providesMessages := request.(MessagesProvider)
	if !providesTranslates && !providesMessages {
		return nil
	}
	translator := &Translator{}
	applyRequestTranslations(request, translator)
	return translator.Compile()
}

// applyRequestTranslations 应用请求翻译配置
func applyRequestTranslations(request any, translator *Translator) {
	if provider, ok := request.(TranslatesProvider); ok {
		translator.Translates = provider.Translates()
	}
	if provider, ok := request.(MessagesProvider); ok {
		translator.Messages = provider.Messages()
	}
}
