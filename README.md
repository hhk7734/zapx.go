## zapx.go

### Context

```go
package main

import (
	"context"

	"github.com/hhk7734/zapx.go"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)

	A(ctx)
}

func A(ctx context.Context) {
	logger := zapx.Ctx(ctx)
	logger.Info("A")

	B(zapx.WithFields(ctx, zap.String("A", "a")))
}

func B(ctx context.Context) {
	logger := zapx.Ctx(ctx)
	logger.Info("B")

	C(zapx.With(ctx, logger.WithOptions(zap.Fields(zap.String("B", "b")))))
}

func C(ctx context.Context) {
	logger := zapx.Ctx(ctx)
	logger.Info("C")
}
```

```json
{"level":"info","ts":1257894000,"caller":"main.go:23","msg":"A"}
{"level":"info","ts":1257894000,"caller":"main.go:30","msg":"B","A":"a"}
{"level":"info","ts":1257894000,"caller":"main.go:37","msg":"C","A":"a","B":"b"}
```

### Dict

```go
package main

import (
	"github.com/hhk7734/zapx.go"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	logger.Info("test",
		zapx.Dict("foo",
			zap.String("bar", "baz")))
}
```

```json
{"level":"info","ts":1257894000,"caller":"main.go:12","msg":"test","foo":{"bar":"baz"}}
```