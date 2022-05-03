package test

import (
	"context"
	"testing"

	pg "github.com/conduitio/conduit-connector-postgres"
	"github.com/conduitio/conduit-connector-postgres/destination"
	"github.com/conduitio/conduit-connector-postgres/source"
	sdk "github.com/conduitio/conduit-connector-sdk"
)

var tableUnderTest string

func init() {
	ctx := context.Background()
	conn := ConnectSimple(ctx, &testing.T{}, RepmgrConnString)

	// todo: this is gross, but how else do I acccess the table made?
	tableUnderTest = SetupTestTable(ctx, &testing.T{}, conn)
}

func TestAcceptance(t *testing.T) {
	sdk.AcceptanceTest(t, sdk.ConfigurableAcceptanceTestDriver{
		Config: sdk.ConfigurableAcceptanceTestDriverConfig{
			Connector: sdk.Connector{ // Note that this variable should rather be created globally in `connector.go`
				NewSpecification: pg.Specification,
				NewSource:        source.NewSource,
				NewDestination:   destination.NewDestination,
			},
			SourceConfig: map[string]string{
				"table":   tableUnderTest,
				"url":     RegularConnString,
				"columns": "id,key",
			},
			DestinationConfig: map[string]string{
				"url":   RegularConnString,
				"table": "records",
			},
			BeforeTest: func(t *testing.T) {
				t.Logf("table under test: %v", tableUnderTest)
			},
		},
	})
}
