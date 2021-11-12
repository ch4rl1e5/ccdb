package main

type ExampleService interface {
	GetStream()
}

type impl struct {
}

func NewExampleService() ExampleService {
	return &impl{}
}

func (e *impl) GetStream() {
	return
}
