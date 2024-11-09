package utils

type MinMaxGeneric[T any] struct{ Min, Max T }

type ResultErr[T any] struct {
	Data  T
	Error error
}

type Entry2[T any, V any] struct {
	A T
	B V
}
