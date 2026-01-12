package iter

import (
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/option"
)

func Last[U any](it Seq[U]) Option[U] {
	return MonadReduce(MonadMap(it, option.Of[U]), function.SK, option.None[U]())
}
