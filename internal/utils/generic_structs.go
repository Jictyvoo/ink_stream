package utils

type MinMaxGeneric[T any] struct{ Min, Max T }

type ResultErr[T any] struct {
	Data  T
	Error error
}
