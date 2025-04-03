package models

var registeredModels []interface{}

// RegisterModel adds a model to the registry
func RegisterModel(model interface{}) {
	registeredModels = append(registeredModels, model)
}

// GetRegisteredModels returns all registered models
func GetRegisteredModels() []interface{} {
	return registeredModels
}
