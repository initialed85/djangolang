package server

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/initialed85/djangolang/pkg/djangolang"
	"github.com/initialed85/djangolang/pkg/helpers"
	"golang.org/x/exp/maps"
)

var (
	logger = helpers.GetLogger("server")
)

func Run(ctx context.Context) error {
	port, err := helpers.GetPort()
	if err != nil {
		return err
	}

	db, err := helpers.GetDBFromEnvironment(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = db.Close()
	}()

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	selectHandlerByEndpointName, err := djangolang.GetSelectHandlerByEndpointName(db)
	if err != nil {
		return err
	}

	endpointNames := maps.Keys(selectHandlerByEndpointName)
	slices.Sort(endpointNames)

	for _, endpointName := range endpointNames {
		selectHandler := selectHandlerByEndpointName[endpointName]
		pathName := fmt.Sprintf("/%v", endpointName)
		logger.Printf("registering %v", pathName)
		r.Get(pathName, selectHandler)
	}

	http.ListenAndServe(
		fmt.Sprintf(":%v", port),
		r,
	)

	return nil
}
