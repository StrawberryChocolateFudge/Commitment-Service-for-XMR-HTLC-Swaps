package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gabstv/httpdigest"

	"gitlab.com/moneropay/go-monero/walletrpc"
)

func getClient() *walletrpc.Client {
	moneroauthuser := os.Getenv("moneroauthuser")
	moneroauthpass := os.Getenv("moneroauthpass")

	return walletrpc.New(walletrpc.Config{
		Address: "http://127.0.0.1:18083/json_rpc",
		Client: &http.Client{
			Transport: httpdigest.New(moneroauthuser, moneroauthpass),
		},
	})
}

func validateAddress(addr string) (bool, error) {
	client := getClient()
	res, err := client.ValidateAddress(context.Background(), &walletrpc.ValidateAddressRequest{
		Address: addr,
	})

	if err != nil {
		return false, fmt.Errorf("addr: %x: %v", addr, err)
	}

	return res.Valid, nil
}

func checkTxKey(Txid string, TxKey string, Address string) (*walletrpc.CheckTxKeyResponse, error) {
	client := getClient()

	res, err := client.CheckTxKey(context.Background(), &walletrpc.CheckTxKeyRequest{
		Txid:    Txid,
		TxKey:   TxKey,
		Address: Address,
	})

	if err != nil {
		return nil, fmt.Errorf("Error when looking up Txid: %v", err)
	}

	return res, nil
}

func fromViewKey(restoreHeight uint64, viewKey string, address string, filename string, password string, save_current bool) (*walletrpc.GenerateFromKeysResponse, error) {
	client := getClient()

	res, err := client.GenerateFromKeys(context.Background(), &walletrpc.GenerateFromKeysRequest{
		RestoreHeight:   restoreHeight,
		Filename:        filename,
		Address:         address,
		ViewKey:         viewKey,
		Password:        password,
		AutosaveCurrent: save_current,
	})
	if err != nil {
		return nil, fmt.Errorf("Error when creating view only wallet V%", err)
	}
	return res, nil
}

func handleViewKeyResponse(info string) {
	// TODO: handle the info to see if the wallet was generated
}

func getTransfer(txid string) (*walletrpc.GetTransferByTxidResponse, error) {
	client := getClient()
	res, err := client.GetTransferByTxid(context.Background(), &walletrpc.GetTransferByTxidRequest{
		Txid: txid,
	})

	if err != nil {
		return nil, fmt.Errorf("Error when getting transfer v%", err)
	}

	return res, nil
}

func handleTransferDetails() {
	//TODO:handle how the transfer was done
}
