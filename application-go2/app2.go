/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"encoding/json"
	"strings"
	"bufio"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

type Asset struct {
	ID             string `json:"ID"`
	Model          string `json:"Model"`
	Price		int    `json:"Price"`
	Color          string `json:"Color"`
	Fuel           string    `json:"Fuel"`
}

//Realizzo un wallet di Token NFT
type Wallet struct{
	Owner string	`json:"Owner"`
	NFT	map[string]Asset	`json:"NFT"`
	}

func main() {
	//Questa roba è tutta magia
	log.Println("============ AVVIO app per Org2MSP (Concessionaria) ============")

	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}

	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}

	if !wallet.Exists("appUser") {
		err = populateWallet(wallet)
		if err != nil {
			log.Fatalf("Failed to populate wallet contents: %v", err)
		}
	}

	ccpPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org2.example.com",
		"connection-org2.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()
//Ottengo il riferimento al canale che mi interessa
	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}
//Ottengo il riferimento al chaincode che mi interessa
	contract1 := network.GetContract("basic-1")
	contract2 := network.GetContract("basic-2")
	reader := bufio.NewReader(os.Stdin)
	log.Println("--> Controllo se la fabbrica ha inizializzato")
	result , err := contract1.EvaluateTransaction("GetAllWallets")
	var wallets []Wallet
	err = json.Unmarshal(result, &wallets)
	if err != nil {
		log.Println("--> Inizializzare i veicoli della fabbrica")
			return 
		
	} else {
		log.Println("--> Fabbrica già inizializzata")
	}
	
	log.Println("--> Controllo se la concessionaria è stato inizializzata")
	result , err = contract2.EvaluateTransaction("GetAllWallets")
	err = json.Unmarshal(result, &wallets)
	if err != nil {
		log.Println("--> Inizializzo la concessionaria")
		_ , err = contract2.SubmitTransaction("InitLedger")
		if err != nil {
			log.Fatalf("Failed to Submit transaction: %v", err)
			return 
		}
	} else {
		log.Println("--> concessionaria già inizializzata")
	}
	
	
	//Inizio
	for {
		fmt.Print("Cosa vuoi fare?\n1. Elenco auto in deposito\n2. Elenco auto in fabbrica\n3. Vendita auto\n4. Ordina auto da fabbrica\n5. Exit:\n ")
		op, _ := reader.ReadString('\n')
		op = strings.Replace(op, "\n", "", -1)
		switch op {
			case "1":
				log.Println("--> Submit Transaction: Elenco auto in deposito")
				result, err := contract2.EvaluateTransaction("GetWallet")
				if err != nil {
					log.Println("Failed to Submit transaction: %v", err)
					break
				}
				log.Println(string(result))
				fmt.Print("\n")
				
			case "2":
				log.Println("--> Evaluate Transaction: Elenco auto in fabbrica")
				result, err := contract1.EvaluateTransaction("GetWallet")
				if err != nil {
					log.Println("Failed to evaluate transaction: %v", err)
					break
				}
				log.Println(string(result))
				fmt.Print("\n")
				
			case "3":
				log.Println("--> Evaluate Transaction: Vendita auto")
				fmt.Print("Inserisci ID: ")
				id, _ := reader.ReadString('\n')
				id = strings.Replace(id, "\n", "", -1)
				_, err := contract2.SubmitTransaction("DeleteAsset", id)
				if err != nil {
					log.Println("Failed to evaluate transaction: %v\n", err)
					break
				}
				
				fmt.Println("AUTO VENDUTA")
				fmt.Print("\n")
				
			case "4":
				log.Println("--> Evaluate Transaction: Ordina auto da fabbrica")
				fmt.Print("Elenco auto disponibili nel deposito in fabbrica:")
				fmt.Print("\n")
				result, err := contract1.EvaluateTransaction("GetWallet")
				if err != nil {
					log.Println("Failed to evaluate transaction: %v", err)
					break
				}
				fmt.Println(string(result))
				fmt.Print("Inserisci ID della vettura da ordinare: ")
				id, _ := reader.ReadString('\n')
				id = strings.Replace(id, "\n", "", -1)
				//fmt.Println("ID PRESO: ",id)
				result, err = contract1.SubmitTransaction("DeleteAsset", id)
				if err != nil {
					log.Println("Failed to evaluate transaction: %v\n", err)
					break
				}
				m := make(map[string]interface{})
				err3 := json.Unmarshal(result, &m)
				if err != nil {
				    log.Fatal(err3)
				}
				i:=fmt.Sprint(m["ID"])
				c:=fmt.Sprint(m["Color"])
				f:=fmt.Sprint(m["Fuel"])
				p:=fmt.Sprint(m["Price"])
				mo:=fmt.Sprint(m["Model"])
				fmt.Print("\n")
				result2, err2 := contract2.SubmitTransaction("CreateAsset", i, c, f, p, mo)
				if err2 != nil {
					log.Println("Failed to evaluate transaction: %v\n", err2)
					break
				}
				log.Println("AUTO ORDINATA DALLA FABBRICA, SARA' DISPONIBILE ENTRO POCHI GIORNI IN DEPOSITO")
				fmt.Print(string(result2) + "\n")
				
		}//switch
		
		if op == "5" {
			break
		}
	}//for
	
	
	log.Println("============ CHIUSURA APP Org2MSP ============")
	
	log.Println("Cancello la cartella appena creata")
	err = os.RemoveAll("./wallet")
	if err != nil {
		log.Fatalf("ERRORE: %v", err)
	}
	
	
}

func populateWallet(wallet *gateway.Wallet) error {
	log.Println("============ Populating wallet ============")
	credPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org2.example.com",
		"users",
		"User1@org2.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "User1@org2.example.com-cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org2MSP", string(cert), string(key))

	return wallet.Put("appUser", identity)
}




