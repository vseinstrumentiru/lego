package lego

type Validatable interface {
	Validate() error
}
