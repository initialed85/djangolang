package example

import (
	"context"
	"fmt"
	"log"

	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/template/templates" // this is where my templates were generated- yours will be different
)

func Run(ctx context.Context) error {
	db, err := helpers.GetDBFromEnvironment(ctx) // e.g. POSTGRES_DB, POSTGRES_PASSWORD etc
	if err != nil {
		return err
	}
	defer func() {
		_ = db.Close()
	}()

	events, err := templates.SelectEvents(
		ctx,
		db,
		templates.EventTransformedColumns, // you'll mostly want this (all columns + any relevant transforms, e.g. PostGIS Point -> GeoJSON map[string]any)
		helpers.Ptr(fmt.Sprintf("%v DESC", templates.EventStartTimestampColumn)), // order by
		helpers.Ptr(10), // limit
		helpers.Ptr(10), // ofset
		fmt.Sprintf("%v > 1000", templates.EventIDColumn),                // implicitly AND'd together
		fmt.Sprintf("%v < 2000", templates.EventIDColumn),                // implicitly AND'd together
		fmt.Sprintf("%v IN (1, 2)", templates.EventSourceCameraIDColumn), // implicitly AND'd together
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, event := range events {
		log.Printf("%#+v", event)
	}

	return nil
}
