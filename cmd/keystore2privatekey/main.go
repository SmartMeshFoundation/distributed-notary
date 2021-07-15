package main

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
)

//Version version of this build
var Version = "v0.1"

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "private-key-file",
			Usage: "the private-key file name(full name).",
			Value: "UTC--xxxx-xx-xxxxx-xx-xx.xxxxxxxxxx--xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		},
		cli.StringFlag{
			Name:  "password",
			Usage: "the password to unlock the private-key file.",
			Value: "***",
		},
	}
	app.Action = StartMain
	app.Name = "Keystore2privatekey"
	app.Version = Version
	err := app.Run(os.Args)
	if err != nil {
		os.Exit(-1)
	}
}

func StartMain(ctx *cli.Context) {
	privateKeyFile := ctx.String("private-key-file")
	passwd := ctx.String("password")
	privKey, address, err := KeystoreToPrivateKey(privateKeyFile, passwd)
	if err != nil {
		fmt.Printf("quit with err : %s\n", err.Error())
		os.Exit(-1)
	}
	fmt.Println("**********************************************************************")
	fmt.Printf("privKey:[%s]\naddress:[%s]\n", privKey, address)
	fmt.Println("**********************************************************************")
}

func KeystoreToPrivateKey(privateKeyFile, password string) (string, string, error) {
	keyjson, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		fmt.Println("read keyjson file failed：", err)
	}
	unlockedKey, err := keystore.DecryptKey(keyjson, password)
	if err != nil {

		return "", "", err

	}
	privKey := hex.EncodeToString(unlockedKey.PrivateKey.D.Bytes())
	addr := crypto.PubkeyToAddress(unlockedKey.PrivateKey.PublicKey)
	return privKey, addr.String(), nil
}
