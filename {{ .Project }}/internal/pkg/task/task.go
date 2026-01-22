package task

import (
	"context"

	"{{.Computed.common_module_final}}/log"
	"{{.Computed.common_module_final}}/worker"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"

	"{{ .Computed.module_name_final }}/internal/conf"
)

// ProviderSet is task providers.
var ProviderSet = wire.NewSet(New)

// New initializes the task worker from config.
func New(c *conf.Bootstrap) (w *worker.Worker, err error) {
	w = worker.New(
		worker.WithRedisURI(c.Redis.Dsn),
		worker.WithGroup(c.Name),
		worker.WithHandler(func(ctx context.Context, p worker.Payload) error {
			return process(task{
				ctx:     ctx,
				c:       c,
				payload: p,
			})
		}),
	)
	if w.Error != nil {
		log.Error(w.Error)
		err = errors.New("initialize worker failed")
		return
	}

	for id, item := range c.Task.Cron {
		err = w.Cron(
			context.Background(),
			worker.WithRunUUID(id),
			worker.WithRunGroup(item.Name),
			worker.WithRunExpr(item.Expr),
			worker.WithRunTimeout(int(item.Timeout)),
			worker.WithRunMaxRetry(int(item.Retry)),
		)
		if err != nil {
			log.Error(err)
			err = errors.New("initialize worker failed")
			return
		}
	}

	log.Info("initialize worker success")
	return
}

type task struct {
	ctx     context.Context
	c       *conf.Bootstrap
	payload worker.Payload
}

func process(t task) (err error) {
	tr := otel.Tracer("task")
	ctx, span := tr.Start(t.ctx, "Task")
	defer span.End()

	// Use task group to match tasks instead of UID
	// This allows for better organization and reusability
	// Match against the group names defined in config
	switch t.payload.Group {
	case t.c.Task.Group.Every10STask:
		log.WithContext(ctx).Info("every10s task executed: %s", t.payload.Payload)
	case t.c.Task.Group.Every3MinTask:
		log.WithContext(ctx).Info("every3min task executed: %s", t.payload.Payload)
	default:
		log.WithContext(ctx).Warn("unknown task group: %s", t.payload.Group)
	}
	return
}
