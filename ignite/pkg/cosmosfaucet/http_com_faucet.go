package cosmosfaucet

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// package initialization for correct validation of commercionetwork addresses
func init() {
	configTestPrefixes()
}

func configTestPrefixes() {
	AccountAddressPrefix := "did:com:"
	AccountPubKeyPrefix := AccountAddressPrefix + "pub"
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, AccountPubKeyPrefix)
	config.Seal()
}

func (f Faucet) comFaucetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	address, err := sdk.AccAddressFromBech32(vars["addr"])
	if err != nil {
		responseError(w, http.StatusBadRequest, err)
		return
	}

	var coins []sdk.Coin
	coin, err := sdk.ParseCoinNormalized(vars["amount"] + "ucommercio")
	if err != nil {
		responseError(w, http.StatusBadRequest, err)
		return
	}
	coins = append(coins, coin)

	// try performing the transfer
	if err := f.Transfer(r.Context(), address.String(), coins); err != nil {
		if err == context.Canceled {
			return
		}
		responseError(w, http.StatusInternalServerError, err)
	} else {
		responseSuccess(w)
	}
}
