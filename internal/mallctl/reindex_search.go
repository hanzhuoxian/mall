package mallctl

import (
	"context"

	marketmodel "github.com/hanzhuoxian/mall/internal/marketplace/model"
	"github.com/hanzhuoxian/mall/internal/marketplace/search"
	marketstore "github.com/hanzhuoxian/mall/internal/marketplace/store"
	marketmysql "github.com/hanzhuoxian/mall/internal/marketplace/store/mysql"
	"github.com/hanzhuoxian/mall/internal/pkg/options"
	"github.com/hanzhuoxian/mall/pkg/app"
	"github.com/hanzhuoxian/mall/pkg/nflag"
)

type reindexSearchOptions struct {
	MySQL            *options.MySQLOptions
	ElasticsearchURL string
}

func newReindexSearchCmd() *app.Command {
	opts := &reindexSearchOptions{MySQL: options.NewMySQLOptions()}
	opts.MySQL.Database = "mall"
	return app.NewCommand("reindex-search", "Rebuild the Elasticsearch product catalog index",
		app.WithCommandOptions(opts),
		app.WithCommandRunFunc(func(args []string) error { return runReindexSearch(opts) }),
	)
}

func (o *reindexSearchOptions) Flags() nflag.NamedFlagSets {
	var sets nflag.NamedFlagSets
	o.MySQL.AddFlags(sets.FlagSet("mysql"))
	sets.FlagSet("search").StringVar(&o.ElasticsearchURL, "elasticsearch-url", "http://localhost:9200", "Elasticsearch URL")
	return sets
}

func (o *reindexSearchOptions) Validate() []error { return o.MySQL.Validate() }

func runReindexSearch(opts *reindexSearchOptions) error {
	db, err := opts.MySQL.NewClient()
	if err != nil {
		return err
	}
	defer func() {
		sqlDB, e := db.DB()
		if e == nil {
			_ = sqlDB.Close()
		}
	}()
	ds := marketmysql.NewStore(db)
	index := search.New(opts.ElasticsearchURL)
	return index.Reindex(context.Background(), func(ctx context.Context, pageOpts marketmodel.ListOptions) (marketmodel.Page[marketmodel.Product], error) {
		return catalogPage(ctx, ds, pageOpts)
	})
}

func catalogPage(ctx context.Context, ds marketstore.Store, opts marketmodel.ListOptions) (marketmodel.Page[marketmodel.Product], error) {
	return ds.ListProducts(ctx, marketmodel.ProductFilter{Catalog: true}, opts)
}
