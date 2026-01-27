package effect

import (
	"errors"
	"fmt"
	"testing"

	"github.com/IBM/fp-go/v2/context/readerioresult"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/codec"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	"github.com/IBM/fp-go/v2/reader"
)

type (
	InjectionToken string

	InjectionContainer = Effect[InjectionToken, any]
)

const (
	service1 = InjectionToken("service1")
	service2 = InjectionToken("service2")
)

func makeSampleInjectionContainer() InjectionContainer {

	return func(token InjectionToken) ReaderIOResult[any] {
		switch token {
		case service1:
			return readerioresult.Of(any("sevice1"))
		case service2:
			return readerioresult.Of(any("sevice2"))
		default:
			return readerioresult.Left[any](errors.New("dependency not available"))
		}

	}
}

var (
	stringCodec = codec.String()

	getService1 = F.Flow2(
		reader.Read[ReaderIOResult[any]](service1),
		readerioresult.ChainEitherK(F.Flow2(
			stringCodec.Decode,
			validation.ToResult[string],
		)),
	)
)

func handleService1() Effect[string, string] {
	return func(ctx string) ReaderIOResult[string] {
		return readerioresult.Of(fmt.Sprintf("Service1: %s", ctx))
	}
}

func TestDependencyLookup(t *testing.T) {

	container := makeSampleInjectionContainer()

	handle := F.Pipe1(
		handleService1(),
		LocalReaderIOResultK[string](getService1),
	)

	res := handle(container)(t.Context())()

	fmt.Println(res)

}
