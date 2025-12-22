{{ if .Computed.enable_db_final }}
{{ if eq .Computed.orm_type_final "none" }}
package data

import (
{{- if .Computed.enable_redis_final }}
	"context"
{{- end }}
	"database/sql"
	"errors"
{{- if .Computed.enable_redis_final }}
	"net/url"
{{- end }}
	"strings"
{{- if .Computed.enable_redis_final }}
	"time"
{{- end }}

	"github.com/go-cinch/common/log"
{{- if .Computed.enable_redis_final }}
	"github.com/go-cinch/common/utils"
{{- end }}
{{- if eq .Computed.db_type_final "postgres" }}
	_ "github.com/lib/pq"
{{- else if eq .Computed.db_type_final "mysql" }}
	_ "github.com/go-sql-driver/mysql"
{{- end }}
{{- if .Computed.enable_redis_final }}
	"github.com/redis/go-redis/v9"
{{- end }}

	"{{.Computed.module_name_final}}/internal/conf"
)

// Data wraps a plain sql.DB connection.
type Data struct {
	DB *sql.DB
}

// NewData initializes a database connection via database/sql.
func NewData(c *conf.Bootstrap) (*Data, func(), error) {
	if c == nil || c.Db == nil {
		err := errors.New("db config is required")
		log.WithError(err).Error("initialize data failed")
		return nil, nil, err
	}

	dsn := strings.TrimSpace(c.Db.Dsn)
	if dsn == "" {
		err := errors.New("db DSN is required")
		log.WithError(err).Error("initialize data failed")
		return nil, nil, err
	}

	const driver = "{{ .Computed.db_type_final }}"

	db, err := sql.Open(driver, dsn)
	if err != nil {
		log.WithError(err).Error("open database failed")
		return nil, nil, err
	}

	if err := db.Ping(); err != nil {
		log.WithError(err).Error("ping database failed")
		return nil, nil, err
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)

	cleanup := func() {
		if err := db.Close(); err != nil {
			log.WithError(err).Error("close database failed")
		}
	}

	log.Info("initialize database success, driver: %s", driver)

	return &Data{DB: db}, cleanup, nil
}
{{- if .Computed.enable_redis_final }}

// NewRedis initializes Redis client from config.
func NewRedis(c *conf.Bootstrap) (redis.UniversalClient, error) {
	return newRedis(c)
}
{{- end }}
{{- else if eq .Computed.orm_type_final "gorm" }}
package data

import (
	"context"
	"errors"
{{- if .Computed.enable_redis_final }}
	"net/url"
{{- end }}
	"strconv"
	"strings"
	"time"

	"github.com/go-cinch/common/id"
	"github.com/go-cinch/common/log"
{{- if .Computed.enable_redis_final }}
	"github.com/go-cinch/common/utils"
{{- end }}
	glog "github.com/go-cinch/common/plugins/gorm/log"
	"github.com/go-cinch/common/plugins/gorm/tenant/v2"
{{- if .Computed.enable_redis_final }}
	"github.com/redis/go-redis/v9"
{{- end }}
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

{{- if .Computed.enable_biz_tx_final }}
	"{{.Computed.module_name_final}}/internal/biz"
{{- end }}
	"{{.Computed.module_name_final}}/internal/conf"
	"{{.Computed.module_name_final}}/internal/db"
)

// Data wraps all data sources used by the service.
type Data struct {
	Tenant    *tenant.Tenant
	sonyflake *id.Sonyflake
}

// NewData initializes the configured database connection via tenant v2.
func NewData(c *conf.Bootstrap) (*Data, func(), error) {
	gormTenant, err := NewDB(c)
	if err != nil {
		return nil, nil, err
	}

	sonyflake, err := NewSonyflake(c)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		log.Info("closing database connections")
	}

	return &Data{
		Tenant:    gormTenant,
		sonyflake: sonyflake,
	}, cleanup, nil
}

// NewDB initializes tenant-aware database connection using tenant v2 package.
// Supports both MySQL and PostgreSQL via internal driver detection.
func NewDB(c *conf.Bootstrap) (*tenant.Tenant, error) {
	if c == nil || c.Db == nil {
		err := errors.New("db config is required")
		log.WithError(err).Error("initialize db failed")
		return nil, err
	}

	dbConf := c.Db
	driver := strings.ToLower(strings.TrimSpace(dbConf.Driver))
	dsn := strings.TrimSpace(dbConf.Dsn)
	if driver == "" {
		err := errors.New("db driver is required")
		log.WithError(err).Error("initialize db failed")
		return nil, err
	}
	if dsn == "" {
		err := errors.New("db DSN is required")
		log.WithError(err).Error("initialize db failed")
		return nil, err
	}

	level := log.NewLevel(c.Log.Level)
	// force to warn level when show sql is false
	if level > log.WarnLevel && !c.Log.ShowSQL {
		level = log.WarnLevel
	}

	ops := []func(*tenant.Options){
		tenant.WithDriver(driver),
		tenant.WithDSN("", dsn), // Empty string for default tenant
		tenant.WithSQLFile(db.SQLFiles),
		tenant.WithSQLRoot(db.SQLRoot),
		tenant.WithSkipMigrate(!dbConf.Migrate), // Skip migration if Migrate is false
		tenant.WithConfig(&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			QueryFields: true,
			Logger: glog.New(
				glog.WithColorful(false),
				glog.WithSlow(200),
				glog.WithLevel(level),
			),
		}),
		tenant.WithMaxIdle(10),
		tenant.WithMaxOpen(100),
	}

	gormTenant, err := tenant.New(ops...)
	if err != nil {
		log.WithError(err).Error("create tenant failed")
		return nil, err
	}

	// Always call Migrate() to initialize the database connection
	// WithSkipMigrate controls whether SQL migrations are actually executed
	if err := gormTenant.Migrate(); err != nil {
		log.WithError(err).Error("migrate tenant failed")
		return nil, err
	}

	log.Info("initialize db success, driver: %s", driver)
	return gormTenant, nil
}

type contextTxKey struct{}

// Tx is transaction wrapper.
func (d *Data) Tx(ctx context.Context, handler func(ctx context.Context) error) error {
	return d.Tenant.DB(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return handler(ctx)
	})
}

// DB returns a tenant-aware GORM DB instance from context.
// If a transaction is present in the context, it returns the transaction DB.
func (d *Data) DB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return d.Tenant.DB(ctx)
}

{{- if .Computed.enable_biz_tx_final }}

// NewTransaction creates a new Transaction from Data.
func NewTransaction(d *Data) biz.Transaction {
	return d
}
{{- end }}

// ID generates a unique distributed ID using Sonyflake.
func (d *Data) ID(ctx context.Context) uint64 {
	return d.sonyflake.ID(ctx)
}

// NewSonyflake initializes the Sonyflake ID generator.
func NewSonyflake(c *conf.Bootstrap) (*id.Sonyflake, error) {
	machineID, _ := strconv.ParseUint(c.Server.MachineId, 10, 16)
	sf := id.NewSonyflake(
		id.WithSonyflakeMachineID(uint16(machineID)),
		id.WithSonyflakeStartTime(time.Date({{ now | date "2006" }}, 1, 1, 0, 0, 0, 0, time.UTC)),
	)
	if sf.Error != nil {
		log.WithError(sf.Error).Error("initialize sonyflake failed")
		return nil, errors.New("initialize sonyflake failed")
	}
	log.
		WithField("machine.id", machineID).
		Info("initialize sonyflake success")
	return sf, nil
}
{{- if .Computed.enable_redis_final }}

// NewRedis initializes Redis client from config.
func NewRedis(c *conf.Bootstrap) (redis.UniversalClient, error) {
	return newRedis(c)
}
{{- end }}
{{- end }}
{{- else }}
package data

import (
{{- if .Computed.enable_redis_final }}
	"context"
	"errors"
	"net/url"
	"time"

{{- end }}
	"github.com/go-cinch/common/log"
{{- if .Computed.enable_redis_final }}
	"github.com/go-cinch/common/utils"
	"github.com/redis/go-redis/v9"
{{- end }}

	"{{.Computed.module_name_final}}/internal/conf"
)

// Data represents mock data provider when database is disabled.
type Data struct {
	// Add your mock data fields here.
}

// NewData creates a new mock Data instance.
func NewData(c *conf.Bootstrap) (*Data, func(), error) {
	log.Info("initializing mock data provider (no database)")

	d := &Data{}

	cleanup := func() {
		log.Info("closing mock data provider")
	}

	return d, cleanup, nil
}
{{- if .Computed.enable_redis_final }}

// NewRedis initializes Redis client from config.
func NewRedis(c *conf.Bootstrap) (redis.UniversalClient, error) {
	return newRedis(c)
}
{{- end }}
{{- end }}

{{ if .Computed.enable_redis_final }}
// newRedis is the shared Redis initialization logic.
func newRedis(c *conf.Bootstrap) (client redis.UniversalClient, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var u *url.URL
	u, err = url.Parse(c.Redis.Dsn)
	if err != nil {
		log.Error(err)
		err = errors.New("initialize redis failed")
		return
	}
	u.User = url.UserPassword(u.User.Username(), "***")
	showDsn, _ := url.PathUnescape(u.String())
	client, err = utils.ParseRedisURI(c.Redis.Dsn)
	if err != nil {
		log.Error(err)
		err = errors.New("initialize redis failed")
		return
	}
	err = client.Ping(ctx).Err()
	if err != nil {
		log.Error(err)
		err = errors.New("initialize redis failed")
		return
	}
	log.
		WithField("redis.dsn", showDsn).
		Info("initialize redis success")
	return
}
{{ end }}
