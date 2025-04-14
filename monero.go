package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/StrawberryChocolateFudge/go-monero/walletrpc"
)

func getClient() *walletrpc.Client {
	//	moneroauthuser := os.Getenv("moneroauthuser")
	//	moneroauthpass := os.Getenv("moneroauthpass")

	return walletrpc.New(walletrpc.Config{
		Address: "http://localhost:18083/json_rpc",
		Client:  &http.Client{
			// Transport: httpdigest.New("kernal", "s3cure"), //TODO: temporary test pass
		},
	})
}

func getBalance() {
	client := getClient()
	resp, err := client.GetBalance(context.Background(), &walletrpc.GetBalanceRequest{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
	fmt.Println("Total balance:", walletrpc.XMRToDecimal(resp.Balance))
	fmt.Println("Unlocked balance:", walletrpc.XMRToDecimal(resp.UnlockedBalance))
}

func getHeight() (*walletrpc.GetHeightResponse, error) {
	client := getClient()
	res, err := client.GetHeight(context.Background())
	if err != nil {
		return nil, fmt.Errorf("addr: %v", err)
	}
	return res, nil
}

func validateAddress(addr string) (bool, error) {
	client := getClient()
	res, err := client.ValidateAddress(context.Background(), &walletrpc.ValidateAddressRequest{
		Address:        addr,
		AnyNetType:     false,
		AllowOpenalias: true,
	})

	if err != nil {
		return false, fmt.Errorf("addr: %s: %v", addr, err)
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
		return nil, fmt.Errorf("v%", err)
	}
	return res, nil
}

func openWallet(filename string, password string) error {
	client := getClient()
	return client.OpenWallet(context.Background(), &walletrpc.OpenWalletRequest{
		Filename: filename,
		Password: password,
	})
}

func handleViewKeyResponse(info string) {
	// TODO: handle the info to see if the wallet was generated
}

func getIncomingTransfers(accountIndex uint64) (*walletrpc.IncomingTransfersResponse, error) {
	client := getClient()
	res, err := client.IncomingTransfers(context.Background(), &walletrpc.IncomingTransfersRequest{
		TransferType: "all",
		AccountIndex: accountIndex,
	})

	if err != nil {
		return nil, fmt.Errorf("v%", err)
	}

	return res, nil
}

func getTransfer(txid string, accountIndex uint64) (*walletrpc.GetTransferByTxidResponse, error) {
	client := getClient()
	res, err := client.GetTransferByTxid(context.Background(), &walletrpc.GetTransferByTxidRequest{
		Txid:         txid,
		AccountIndex: accountIndex,
	})

	if err != nil {
		return nil, fmt.Errorf("v%", err)
	}

	return res, nil
}

func getAddressIndex(addr string) (*walletrpc.GetAddressIndexResponse, error) {
	client := getClient()
	res, err := client.GetAddressIndex(context.Background(), &walletrpc.GetAddressIndexRequest{
		Address: addr,
	})

	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return res, nil
}

func getAddress() (*walletrpc.GetAddressResponse, error) {
	client := getClient()
	res, err := client.GetAddress(context.Background(), &walletrpc.GetAddressRequest{
		AccountIndex: 0,
	})

	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return res, nil
}

func handleTransferDetails() {
	//TODO:handle how the transfer was done
}

func rescanBlockchain() error {
	client := getClient()
	err := client.RescanBlockchain(context.Background())

	if err != nil {
		return fmt.Errorf("%v", err)
	}
	return err
}
