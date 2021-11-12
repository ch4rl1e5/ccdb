package stream

type Client interface {
}

type implClient struct {
}

func NewClient() Client {
	return &implClient{}
}
