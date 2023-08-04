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
	defer logger.Sync()
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
{"level":"info","ts":1689501775.2534308,"caller":"zapx.go/main.go:22","msg":"A"}
{"level":"info","ts":1689501775.2534587,"caller":"zapx.go/main.go:29","msg":"B","A":"a"}
{"level":"info","ts":1689501775.2534641,"caller":"zapx.go/main.go:36","msg":"C","A":"a","B":"b"}
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
	defer logger.Sync()

	logger.Info("test",
		zapx.Dict("foo",
			zap.String("bar", "baz")))
}
```

```json
{"level":"info","ts":1689501801.7570314,"caller":"zapx.go/main.go:12","msg":"test","foo":{"bar":"baz"}}
```

### GormLogger

```go
package main

import (
	"context"

	"github.com/hhk7734/zapx.go"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string
}

func main() {
	ctx := context.Background()

	cfg := zap.NewProductionConfig()
	cfg.Level.SetLevel(zap.DebugLevel)
	logger, _ := cfg.Build()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	db, _ := gorm.Open(
		sqlite.Open("test.db"),
		&gorm.Config{
			Logger: zapx.DefaultGormLogger(),
		})

	db.Debug().AutoMigrate(&User{})

	ctx = zapx.WithFields(ctx, zap.String("test", "test"))
	db = db.WithContext(ctx)

	db.Debug().Create(&User{Name: "test"})
}
```

```json
{"level":"debug","ts":1689501127.6352503,"caller":"zapx.go/main.go:32","msg":"trace: debug","elapsed":0.000024008,"rows":-1,"sql":"SELECT count(*) FROM sqlite_master WHERE type='table' AND name=\"users\""}
{"level":"debug","ts":1689501127.6410747,"caller":"zapx.go/main.go:32","msg":"trace: debug","elapsed":0.005701533,"rows":0,"sql":"CREATE TABLE `users` (`id` integer,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`name` text,PRIMARY KEY (`id`))"}
{"level":"debug","ts":1689501127.6445901,"caller":"zapx.go/main.go:32","msg":"trace: debug","elapsed":0.003401466,"rows":0,"sql":"CREATE INDEX `idx_users_deleted_at` ON `users`(`deleted_at`)"}

{"level":"debug","ts":1689501127.6482286,"caller":"zapx.go/main.go:37","msg":"trace: debug","test":"test","elapsed":0.00357918,"rows":1,"sql":"INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`name`) VALUES (\"2023-07-16 18:52:07.644\",\"2023-07-16 18:52:07.644\",NULL,\"test\") RETURNING `id`"}
```

If you want to ignore ErrRecordNotFound, set `IgnoreRecordNotFoundError` to true.

```go
	db, _ := gorm.Open(
		sqlite.Open("test.db"),
		&gorm.Config{
			Logger: &zapx.GormLogger{
				Config: zapx.Config{
					SlowThreshold:             200 * time.Millisecond,
					Colorful:                  false, // not support
					IgnoreRecordNotFoundError: true,
					ParameterizedQueries:      false,
					LogLevel:                  zapx.Warn,
				},
			},
		})
```