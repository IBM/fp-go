package readerreaderioresult

import (
	"strconv"
	"sync"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	RES "github.com/IBM/fp-go/v2/context/readerioresult"
	"github.com/IBM/fp-go/v2/function"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

type (
	ConsoleDependency interface {
		Log(msg string) IO[Void]
	}

	Res[A any] = RES.ReaderIOResult[A]

	ConsoleEnv[A any] = ReaderReaderIOResult[ConsoleDependency, A]

	consoleOnArray struct {
		logs []string
		mu   sync.Mutex
	}
)

var (
	logConsole = reader.Curry1(ConsoleDependency.Log)
)

func (c *consoleOnArray) Log(msg string) IO[Void] {
	return func() Void {
		c.mu.Lock()
		defer c.mu.Unlock()

		c.logs = append(c.logs, msg)
		return function.VOID
	}
}

func makeConsoleOnArray() *consoleOnArray {
	return &consoleOnArray{}
}

func TestConsoleEnv(t *testing.T) {
	console := makeConsoleOnArray()

	prg := F.Pipe1(
		Of[ConsoleDependency]("Hello World!"),
		TapReaderIOK(logConsole),
	)

	res := prg(console)(t.Context())()

	assert.Equal(t, result.Of("Hello World!"), res)
	assert.Equal(t, A.Of("Hello World!"), console.logs)
}

func TestConsoleEnvWithLocal(t *testing.T) {
	console := makeConsoleOnArray()

	prg := F.Pipe1(
		Of[ConsoleDependency](42),
		TapReaderIOK(reader.WithLocal(logConsole, strconv.Itoa)),
	)

	res := prg(console)(t.Context())()

	assert.Equal(t, result.Of(42), res)
	assert.Equal(t, A.Of("42"), console.logs)
}
