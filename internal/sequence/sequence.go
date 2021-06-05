package sequence

import (
	"github.com/ch4rl1e5/stream/internal/buffer"
	"net/http"
	"reflect"
	"runtime"
)

var registered []FuncSequence

type Sequencies interface {
	BuildSequencies(h http.Header)
	Run(bufferPool buffer.Pools)
}

type FuncSequence func(bufferPool buffer.Pools)  error

type sequence struct {
	FuncSequence
	key string
	config string
}

type Impl []sequence

func New() Sequencies {
	return &Impl{}
}

func (s *Impl) BuildSequencies(h http.Header) {
	for name, values := range h {
		for _, value := range values {
			err := isRegistered(sequence{key: name, config: value})
			if err != nil {
				panic(err)
			}
			*s = append(*s, sequence{key: name, config: value})
		}
	}
}

func (s *Impl) Run(bufferPool buffer.Pools) {
	for _, funcSeq := range registered {
		for _, seq := range *s {
			if runtime.FuncForPC(reflect.ValueOf(funcSeq).Pointer()).Name() == seq.key {
				err := funcSeq(bufferPool)
				if err != nil {
					return
				}

				go func() {
					for {
						bufferPool.Read()
					}
				}()

				go func() {
					for {
						bufferPool.Write()
					}
				}()
			}
		}
	}
}

func isRegistered(seq sequence) error {
	for _, funcSeq := range registered {
		if reflect.TypeOf(funcSeq).Kind() == reflect.Func &&
			runtime.FuncForPC(reflect.ValueOf(funcSeq).Pointer()).Name() == seq.key {
			return nil
		}
	}
	return &ErrSequenceNotRegistered{SequenceKey: seq.key}
}

func RegisterSequence(seq ...FuncSequence) {
	registered = append(registered, seq...)
}