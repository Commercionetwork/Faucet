package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ignite/cli/ignite/pkg/cosmosfaucet"

	"github.com/ignite/cli/ignite/pkg/chaincmd"
	chaincmdrunner "github.com/ignite/cli/ignite/pkg/chaincmd/runner"
	"github.com/spf13/viper"
)

type Config struct {
	ChainID         string `mapstructure:"ChainID"`
	LcdNode         string `mapstructure:"LcdNode"`
	MnemonicFaucet  string `mapstructure:"MnemonicFaucet"`
	AccountFaucet   string `mapstructure:"AccountFaucet"`
	AccountCoinType string `mapstructure:"AccountCoinType"`
	RpcNode         string `mapstructure:"RpcNode"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, err := LoadConfig(".")
	if err != nil {
		panic(err)
	}

	chainID := "commercio-testnet11k"
	lcdNode := "https://lcd-testnet.commercio.network"
	accountFaucet := "faucettestnet"
	accountCoinType := "118"
	rpcNode := "https://rpc-testnet.commercio.network:443"

	mnemonicFaucet := config.MnemonicFaucet

	if mnemonicFaucet == "" {
		panic(fmt.Errorf("Faucet mnemonic can't be empty!"))
	}

	if config.ChainID != "" {
		chainID = config.ChainID
	}
	if config.LcdNode != "" {
		lcdNode = config.LcdNode
	}
	if config.AccountFaucet != "" {
		accountFaucet = config.AccountFaucet
	}
	if config.AccountCoinType != "" {
		accountCoinType = config.AccountCoinType
	}
	if config.RpcNode != "" {
		rpcNode = config.RpcNode
	}

	commandPath := "commercionetworkd"

	chainCommandOptions := []chaincmd.Option{
		chaincmd.WithChainID(chainID),
		chaincmd.WithFees("10000ucommercio"),
		chaincmd.WithKeyringBackend("test")
		//chaincmd.WithNodeAddress("tcp://64.225.78.169:26657"),
		chaincmd.WithNodeAddress(rpcNode),
	}

	cmd := chaincmd.New(commandPath, chainCommandOptions...)

	runner, err := chaincmdrunner.New(ctx, cmd)
	if err != nil {
		panic(err)
	}

	comDenom := "ucommercio"

	faucet, err := cosmosfaucet.New(ctx, runner,
		cosmosfaucet.Account(accountFaucet, mnemonicFaucet, accountCoinType),
		cosmosfaucet.Coin(1000, 2000, comDenom),
		cosmosfaucet.ChainID(chainID),
		cosmosfaucet.OpenAPI(lcdNode),
	)
	if err != nil {
		panic(err)
	}

	http.ListenAndServe(":8181", faucet)
}
