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
	log.Println("============ AVVIO app per Org1MSP (Fabbrica) ============")

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
		"org1.example.com",
		"connection-org1.yaml",
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
	contract := network.GetContract("basic-1")
	contract2 := network.GetContract("basic-2")
	
	reader := bufio.NewReader(os.Stdin)
	log.Println("--> Controllo se il ledger è stato inizializzato")
	result , err := contract.EvaluateTransaction("GetAllWallets")
	var wallets []Wallet
	err = json.Unmarshal(result, &wallets)
	if err != nil {
		log.Println("--> Inizializzo il ledger")
		_ , err = contract.SubmitTransaction("InitLedger")
		if err != nil {
			log.Fatalf("Failed to Submit transaction: %v", err)
			return 
		}
	} else {
		log.Println("--> Ledger già inizializzato")
	}
	
	
	//Inizio
	for {
		fmt.Print("Cosa vuoi fare?\n1. InitLedger\n2. GetAllAssets\n3. ReadAsset\n4. UpdateAsset\n5. DeleteAsset\n6. AssetExists\n7. TransferAsset\n8. CreateAsset\n9. Exit:\n ")
		op, _ := reader.ReadString('\n')
		op = strings.Replace(op, "\n", "", -1)
		switch op {
			case "1":
				log.Println("--> Submit Transaction: InitLedger")
				result, err := contract.SubmitTransaction("InitLedger")
				if err != nil {
					log.Println("Failed to Submit transaction: %v", err)
					break
				}
				log.Println("ESEGUITO")
				fmt.Print(string(result) + "\n")
				
			case "2":
				log.Println("--> Evaluate Transaction: GetAllAssets")
				fmt.Print("\n")
				fmt.Println("                               Elenco asset in deposito fabbrica")
				result, err := contract.EvaluateTransaction("GetWallet")
				if err != nil {
					log.Println("Failed to evaluate transaction: %v", err)
					break
				}
				fmt.Println(string(result))
				fmt.Print("\n")
				fmt.Println("                          Elenco asset in deposito concessionaria")
				result2, err2 := contract2.EvaluateTransaction("GetWallet")
				if err2 != nil {
					log.Println("Failed to evaluate transaction: %v", err2)
					break
				}
				fmt.Println(string(result2))
				fmt.Print("\n")
				
			case "3":
				log.Println("--> Evaluate Transaction: ReadAsset")
				fmt.Print("Inserisci ID: ")
				id, _ := reader.ReadString('\n')
				id = strings.Replace(id, "\n", "", -1)
				result, err := contract.EvaluateTransaction("ReadAsset", id)
				if err != nil {
					log.Println("Failed to evaluate transaction: %v\n", err)
					break
				}
				log.Println(string(result))
				fmt.Print("\n")
				
			case "4":
				log.Println("--> Evaluate Transaction: UpdateAsset")
				fmt.Print("Inserisci ID: ")
				id, _ := reader.ReadString('\n')
				id = strings.Replace(id, "\n", "", -1)
				fmt.Print("Inserisci modello: ")
				model, _ := reader.ReadString('\n')
				model = strings.Replace(model, "\n", "", -1)
				fmt.Print("Inserisci prezzo: ")
				price, _ := reader.ReadString('\n')
				price = strings.Replace(price, "\n", "", -1)
				fmt.Print("Inserisci colore: ")
				color, _ := reader.ReadString('\n')
				color = strings.Replace(color, "\n", "", -1)
				fmt.Print("Inserisci tipo carburante: ")
				fuel, _ := reader.ReadString('\n')
				fuel = strings.Replace(fuel, "\n", "", -1)
				result, err := contract.SubmitTransaction("UpdateAsset", id, color, model, price, fuel)
				if err != nil {
					log.Println("Failed to evaluate transaction: %v\n", err)
					break
				}
				log.Println("ESEGUITO")
				fmt.Print(string(result) + "\n")
				
			case "5":
				log.Println("--> Evaluate Transaction: DeleteAsset")
				fmt.Print("Inserisci ID: ")
				id, _ := reader.ReadString('\n')
				id = strings.Replace(id, "\n", "", -1)
				result, err := contract.SubmitTransaction("DeleteAsset", id)
				if err != nil {
					log.Println("Failed to evaluate transaction: %v", err)
					break
				}
				log.Println("ESEGUITO")
				fmt.Print(string(result) + "\n")
				
			case "6":
				log.Println("--> Evaluate Transaction: AssetExists")
				fmt.Print("Inserisci ID: ")
				id, _ := reader.ReadString('\n')
				id = strings.Replace(id, "\n", "", -1)
				result, err := contract.EvaluateTransaction("AssetExists", id)
				if err != nil {
					log.Println("Failed to evaluate transaction: %v\n", err)
					break
				}
				result2, err2 := contract2.EvaluateTransaction("AssetExists", id)
				if err2 != nil {
					log.Println("Failed to evaluate transaction: %v\n", err2)
					break
				}
				if string(result2) == "true" {
					fmt.Println("L'asset con ID: " + id + " esiste nel deposito concessionaria")
				} else if string(result2) == "false" {
					fmt.Println("L'asset con ID: " + id + " NON esiste nel deposito concessionaria")
				} else {
					fmt.Println("Errore")
				}
				if string(result) == "true" {
					fmt.Println("L'asset con ID: " + id + " esiste nel deposito fabbrica")
				} else if string(result) == "false" {
					fmt.Println("L'asset con ID: " + id + " NON esiste nel deposito fabbrica")
				} else {
					fmt.Println("Errore")
				}
				fmt.Print("\n")
				
			case "7":
				log.Println("--> Evaluate Transaction: TransferAsset")
				fmt.Print("Inserisci ID: ")
				id, _ := reader.ReadString('\n')
				id = strings.Replace(id, "\n", "", -1)
				result, err := contract.SubmitTransaction("DeleteAsset", id)
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
				log.Println("ESEGUITO")
				fmt.Print(string(result2) + "\n")
				
				
			
			case "8":
				log.Println("--> Evaluate Transaction: CreateAsset")
				fmt.Print("Inserisci ID: ")
				id, _ := reader.ReadString('\n')
				id = strings.Replace(id, "\n", "", -1)
				fmt.Print("Inserisci modello: ")
				model, _ := reader.ReadString('\n')
				model = strings.Replace(model, "\n", "", -1)
				fmt.Print("Inserisci prezzo: ")
				price, _ := reader.ReadString('\n')
				price = strings.Replace(price, "\n", "", -1)
				fmt.Print("Inserisci colore: ")
				color, _ := reader.ReadString('\n')
				color = strings.Replace(color, "\n", "", -1)
				fmt.Print("Inserisci tipo carburante: ")
				fuel, _ := reader.ReadString('\n')
				fuel = strings.Replace(fuel, "\n", "", -1)
				result,_:= contract2.EvaluateTransaction("AssetExists", id)
				if (string(result)=="true") {
					log.Println("Failed to evaluate transaction: ",id," presente nel deposito concessionaria" )
					break
				}
				result2, err2 := contract.SubmitTransaction("CreateAsset", id, color, fuel, price, model)
				if err2 != nil {
					log.Println("Failed to evaluate transaction: %v\n", err2)
					break
				}
				log.Println("ESEGUITO")
				fmt.Print(string(result2) + "\n")
				
					
		}//switch
		
		if op == "9" {
			break
		}
	}//for
	
	
	log.Println("============ CHIUSURA APP Org1MSP ============")
	
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
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "User1@org1.example.com-cert.pem")
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

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	return wallet.Put("appUser", identity)
}





