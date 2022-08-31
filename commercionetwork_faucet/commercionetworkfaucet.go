package main

import (
	"context"
	"net/http"

	"github.com/ignite/cli/ignite/pkg/cosmosfaucet"

	"github.com/ignite/cli/ignite/pkg/chaincmd"
	chaincmdrunner "github.com/ignite/cli/ignite/pkg/chaincmd/runner"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := chaincmd.New("commercionetworkd")

	runner, err := chaincmdrunner.New(ctx, cmd)
	if err != nil {
		panic(err)
	}

	comDenom := "ucommercio"

	faucet, err := cosmosfaucet.New(ctx, runner,
		cosmosfaucet.Account("bob", "special chest leaf section reunion inflict busy blouse inflict kid alcohol hazard embody mosquito green turkey street very lab forest gain disease hollow bomb", "com"),
		cosmosfaucet.Coin(1000, 2000, comDenom),
		cosmosfaucet.ChainID("commercionetwork"),
	)
	if err != nil {
		panic(err)
	}

	http.ListenAndServe(":8181", faucet)
}
