package kerrors

import "testing"

func TestSetDictionaryReturnData(t *testing.T) {
	dictionary := &Dictionary{
		Items: make(map[ErrorCode]ErrorMessage),
	}
	code, msg := ErrorCode("0001"), ErrorMessage("content")
	dictionary.Set(code, msg)
	result := dictionary.Get(code)
	if result != msg {
		t.Errorf("Expected %v, but got %v", msg, result)

	}
}

func TestGetDictionaryReturnNotfound(t *testing.T) {
	dictionary := &Dictionary{
		Items: make(map[ErrorCode]ErrorMessage),
	}
	code, msg := ErrorCode("0001"), ErrorMessage("content")
	invalidCode := ErrorCode("0002")
	dictionary.Set(code, msg)
	result := dictionary.Get(invalidCode)
	if result != NotFoundErrorMessageInDictionary {
		t.Errorf("Expected %v, but got %v", NotFoundErrorMessageInDictionary, result)
		t.Errorf("Expected code %v, but got %v", code.String(), invalidCode.String())

	}
}

func TestHasDictionaryReturnData(t *testing.T) {
	dictionary := &Dictionary{
		Items: make(map[ErrorCode]ErrorMessage),
	}
	code, msg := ErrorCode("0001"), ErrorMessage("content")
	dictionary.Set(code, msg)
	ok := dictionary.Has(code)
	if !ok {
		t.Errorf("Expected true, but got %v", ok)

	}
}

func TestClearDictionaryReturnEmptyItem(t *testing.T) {
	dictionary := &Dictionary{
		Items: make(map[ErrorCode]ErrorMessage),
	}
	code, msg := ErrorCode("0001"), ErrorMessage("content")
	dictionary.Set(code, msg)

	dictionary.Clear()
	if dictionary.Size() != 0 {
		t.Errorf("Expect size is 0, but got size %v", dictionary.Size())
	}
}
