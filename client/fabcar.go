
package main
import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"crypto/md5"
	"encoding/hex"
	"path/filepath"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"database/sql"
    _ "github.com/go-sql-driver/mysql"
)

//Struttura dei dati che andremo a gestire
type Infos struct {
    id int `json:"id"`
    nome string `json:"nome"`
    cognome string `json:"cognome"`
    saldo int `json:"saldo"`

}

func main() {
	os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		fmt.Printf("Failed to create wallet: %s\n", err)
		os.Exit(1)
	}

	if !wallet.Exists("appUser") {
		err = populateWallet(wallet)
		if err != nil {
			fmt.Printf("Failed to populate wallet contents: %s\n", err)
			os.Exit(1)
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
		fmt.Printf("Failed to connect to gateway: %s\n", err)
		os.Exit(1)
	}
	defer gw.Close()

	//Usiamo il canale di default
	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		fmt.Printf("Failed to get network: %s\n", err)
		os.Exit(1)
	}

	//Lasciamo il contratto con il nome di default per semplicita ma potremmo sostituirlo a piacimento
	contract := network.GetContract("fabcar")

	//Dati di default usati per il test
    infoArray := []Infos{
        Infos{-1,"Cataldo", "Baglio", 420},
        Infos{-1,"Giovanni", "Storti",69000},
        Infos{-1,"Giacomo", "Poretti",12500},
        Infos{-1,"Marina","Massironi",25000},
    }

	//Connessione al db mySQL locale
    db, err := sql.Open("mysql", "fabric:password@tcp(127.0.0.1:3306)/tesi")
    if err != nil {
        panic(err.Error())
    }
    defer db.Close()

    fmt.Printf("%s\n\n","================================================================================================================================")
    fmt.Printf("\033[33m%s\033[0m\n","CRUD: Create, Read, Update, Delete")
    fmt.Printf("%s\n\n","================================================================================================================================")
    fmt.Printf("\033[33m%s\033[0m\n","[Create] Inserimento degli hash dei seguenti dati in Fabric")
    fmt.Printf("%s\n","________________________________________________________________________________________________________________________________")
    //Inserimento con un loop
    loop := 0
    for loop < 4{
        //Inserimento nel database mysql
        stmt, _ := db.Prepare("INSERT INTO VeryImportantInfo (nome,cognome,saldo) VALUES (?,?,?);")
        defer stmt.Close()
        res, _ := stmt.Exec(infoArray[loop].nome,infoArray[loop].cognome,strconv.Itoa(infoArray[loop].saldo))

        insertId, _ := res.LastInsertId()
        infoArray[loop].id = int(insertId)
        //Inserimento nella blockchain di Fabric attraverso la funzione del chaincode
	tempHash := GenerateHash(infoArray[loop])
        contract.SubmitTransaction("CreateHash",strconv.Itoa(infoArray[loop].id),tempHash,"Admin")

        fmt.Printf("%s\n",infoArray[loop].nome + " " + infoArray[loop].cognome + " " + strconv.Itoa(infoArray[loop].saldo) + "$ --> " + tempHash)
        loop += 1
        }

	fmt.Printf("%s\n\n","================================================================================================================================")
    fmt.Printf("\033[33m%s\033[0m\n","[Read] Lettura da Fabric:")
    fmt.Printf("%s\n","________________________________________________________________________________________________________________________________")

	//Stampa di tutti i dati presenti su fabric utilizzando la funzione del chaincode
	result, err := contract.EvaluateTransaction("GetAllHashes")
	fmt.Println(string(result))

	//Modifica non autorizzata
	fmt.Printf("%s\n\n","================================================================================================================================")
	fmt.Printf("\033[33m%s\033[0m\n","[Update] Modifica non autorizzata:")
    fmt.Printf("%s\n","________________________________________________________________________________________________________________________________")

    //Select dei dati prima della modifica dal db mySQL
    var info Infos
    err = db.QueryRow("SELECT id, nome, cognome, saldo FROM VeryImportantInfo where id = ?", infoArray[0].id).Scan(&info.id, &info.nome, &info.cognome, &info.saldo)
    oldHash := GenerateHash(info)
    fmt.Println("Dato iniziale : " + info.nome + ", " + info.cognome + ", " + strconv.Itoa(info.saldo) + "$")

	//Update del dato in esame solo sul database mySQL
    stmt, err := db.Prepare("UPDATE VeryImportantInfo SET nome = 'Anlo' where id = ?")
    defer stmt.Close()
    stmt.Exec(strconv.Itoa(infoArray[0].id))

	//Select dei dati dopo la modifica dal db mySQL
    err = db.QueryRow("SELECT id, nome, cognome, saldo FROM VeryImportantInfo where id = ?", infoArray[0].id).Scan(&info.id, &info.nome, &info.cognome, &info.saldo)
    newHash := GenerateHash(info)
    fmt.Println("Dato Finale : " + info.nome + ", " + info.cognome + ", " + strconv.Itoa(info.saldo) + "$")

    //Controllo dell'hash tramite funzione del chaincode
    result, err = contract.EvaluateTransaction("CheckHash", strconv.Itoa(info.id))

    fmt.Printf("%s\n","________________________________________________________________________________________________________________________________")
    fmt.Printf("HASH DB PRIMA = ")
    fmt.Printf("\033[35m%s\033[0m\n", oldHash)
    fmt.Printf("HASH  FABRIC  = ")
    fmt.Printf("\033[35m%s\033[0m\n", string(result))
    fmt.Printf("HASH DB DOPO  = ")
    fmt.Printf("\033[35m%s\033[0m\n", newHash)
    fmt.Printf("%s\n","________________________________________________________________________________________________________________________________")

    if(newHash != string(result)){
        fmt.Printf(string("\033[31m%s\033[0m\n"), "I due hash NON CORRISPONDONO!, i dati sono stati COMPROMESSI")
    }else{
    	fmt.Printf(string("\033[33m%s\033[0m\n"), "I due hash CORRISPONDONO!")
    }

    fmt.Printf("%s\n\n","================================================================================================================================")
    fmt.Printf("\033[33m%s\033[0m\n","[Update] Modifica autorizzata:")
	fmt.Printf("%s\n","________________________________________________________________________________________________________________________________")

	//Select dei dati prima della modifica dal db mySQL
    err = db.QueryRow("SELECT id, nome, cognome, saldo FROM VeryImportantInfo where id = ?", infoArray[0].id).Scan(&info.id, &info.nome, &info.cognome, &info.saldo)
    oldHash = GenerateHash(info)
    fmt.Println("Dato iniziale : " + info.nome + ", " + info.cognome + ", " + strconv.Itoa(info.saldo) + "$")

	//Update del dato in esame solo sul database mySQL
    stmt, err = db.Prepare("UPDATE VeryImportantInfo SET nome = 'Aldo' where id = ?")
    defer stmt.Close()
    stmt.Exec(strconv.Itoa(infoArray[0].id))

	//Select dei dati dopo la modifica dal db mySQL
    db.QueryRow("SELECT id, nome, cognome, saldo FROM VeryImportantInfo where id = ?", infoArray[0].id).Scan(&info.id, &info.nome, &info.cognome, &info.saldo)
    newHash = GenerateHash(info)
    fmt.Println("Dato Finale : " + info.nome + ", " + info.cognome + ", " + strconv.Itoa(info.saldo) + "$")

	//Update del dato anche su Fabric
    result, err = contract.SubmitTransaction("UpdateHash", strconv.Itoa(info.id),GenerateHash(info),"Admin")
	//Controllo dell'hash tramite funzione del chaincode
    result, err = contract.EvaluateTransaction("CheckHash", strconv.Itoa(info.id))

    fmt.Printf("%s\n","________________________________________________________________________________________________________________________________")
    fmt.Printf("HASH DB PRIMA = ")
    fmt.Printf("\033[35m%s\033[0m\n", oldHash)
    fmt.Printf("HASH  FABRIC  = ")
    fmt.Printf("\033[35m%s\033[0m\n", string(result))
    fmt.Printf("HASH DB DOPO  = ")
    fmt.Printf("\033[35m%s\033[0m\n", newHash)
    fmt.Printf("%s\n","________________________________________________________________________________________________________________________________")

    if(newHash != string(result)){
        fmt.Printf(string("\033[31m%s\033[0m\n"), "I due hash NON CORRISPONDONO!, i dati sono stati COMPROMESSI")
    }else{
        fmt.Printf(string("\033[32m%s\033[0m\n"), "I due hash CORRISPONDONO!")
    }

    //Eliminazione
    fmt.Printf("%s\n\n","================================================================================================================================")
    fmt.Printf("\033[33m%s\033[0m\n","[Delete] Eliminazione sia su MySQL sia su Fabric:")
    fmt.Printf("%s\n","________________________________________________________________________________________________________________________________")

	//Stampa di tutti i dati prima dell'eliminazione
    result, err = contract.EvaluateTransaction("GetAllHashes")
    fmt.Println(string(result))

	//Eliminzaione nel database mySQL
    stmt, err = db.Prepare("DELETE FROM VeryImportantInfo WHERE id = ?")
    defer stmt.Close()
    stmt.Exec(strconv.Itoa(infoArray[3].id))

	//Eliminazione del dato su Fabric
    result, err = contract.SubmitTransaction("DeleteHash", strconv.Itoa(infoArray[3].id))
    if err != nil {
        fmt.Printf("Failed to evaluate transaction: %s\n", err)
        os.Exit(1)
    }

    fmt.Printf("%s\n","________________________________________________________________________________________________________________________________")
    result, err = contract.EvaluateTransaction("GetAllHashes")
    fmt.Println(string(result))
    fmt.Printf("%s\n","================================================================================================================================")
}

//Funzione per generare gli hash dei dati
func GenerateHash(info Infos) string{
    stringaTotale := strconv.Itoa(info.id) + info.nome + info.cognome + strconv.Itoa(info.saldo)
    md5 := md5.Sum([]byte(stringaTotale))
    return hex.EncodeToString(md5[:])
}

func populateWallet(wallet *gateway.Wallet) error {
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

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
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
		return errors.New("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	err = wallet.Put("appUser", identity)
	if err != nil {
		return err
	}
	return nil
}
