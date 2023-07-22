package config

type Value interface {
	Bool() (bool, error)
	Int() (int, error)
	Float() (float32, error)
	String() (string, error)
}
