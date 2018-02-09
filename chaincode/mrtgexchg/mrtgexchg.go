/*
 * The sample smart contract for documentation topic:
 * Mortgage processing
 */package main

/* Imports
 * utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"math/rand" //some random number generation for insurance quotes, and fico scores
	"strconv"
	"time" //needed to record transaction history for each time ledger entry is updated

	"github.com/hyperledger/fabric/core/chaincode/shim" // import for Chaincode Interface
	pb "github.com/hyperledger/fabric/protos/peer"      // import for peer response
)

// Defined to implement chaincode interface
type Mrtgexchg struct {
}

//few constants to help generate random numbers
const (
	FicoHigh           = 800.0
	FicoLow            = 600.0
	FicoThreshold      = 650.0
	InsuranceHigh      = 5000.0
	InsuranceLow       = 2500.0
	InsuranceThreshold = 0.0
	AppraisalHigh      = 2000000.0
	AppraisalLow       = 750000.0
	RecordsChaincode   = "recordschaincode"
	RecordsChannel     = "records"
	LendingChaincode   = "lendingchaincode"
	LendingChannel     = "lending"
	BooksChaincode     = "bookschaincode"
	BooksChannel       = "books"
	QueryBooksString   = "queryBooks"
	QueryLendingString = "queryLending"
)

/* -------------------------------------------------------------------------------------------------
Define our struct to store real estates in records Blockchain, start fields upper case for JSON
only Registry can write to the blockchain, all others are readonly
---------------------------------------------------------------------------------------------------*/
type RealEstate struct {
	RealEstateID       string // This one will be our key
	Address            string
	Value              float64
	Details            string // this will contain its status on the exchange
	Owner              string
	TransactionHistory map[string]string
}

/* -------------------------------------------------------------------------------------------------
Define our struct to store customer details  in lending Blockchain, start fields upper case for JSON
Bank, Insurance and Fico can write to this blockchain
 -------------------------------------------------------------------------------------------------*/
type Mortgage struct {
	CustID             string // This one will be our key
	RealEstateID       string //
	LoanAmount         float64
	Fico               float64
	Insurance          float64
	Appraisal          float64           //this we will get from books ledger
	Status             string            //status of the mortgage pending -> Funded -> not Funded
	TransactionHistory map[string]string //to hold details for auditing - includes the function called and timestamp
}

/* -------------------------------------------------------------------------------------------------
// Define our struct to store books (record of the appraisals and titles)  in Blockchain,
start fields upper case for JSON
only Titile and Appraiser can write to this blockchain
 -------------------------------------------------------------------------------------------------*/
type Books struct {
	RealEstateID       string // This one will be our key
	Appraisal          float64
	NewTitleOwner      string
	TitleStatus        bool              //here we will store the results of title search which will be used by bank/lender to close the loan
	TransactionHistory map[string]string //to hold details for auditing - includes the function called and timestamp
}

/* -------------------------------------------------------------------------------------------------
these are utility functions
 -------------------------------------------------------------------------------------------------*/
func getTimeNow() string {
	var formatedTime string
	t := time.Now()
	formatedTime = t.Format(time.RFC1123)
	return formatedTime
}

func random(max, min float64) float64 {
	rand.Seed(time.Now().Unix())
	return rand.Float64()*(max-min) + min
}

func float64ToByte(f float64) []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], math.Float64bits(f))
	return buf[:]
}

/* these functions are to generate randomly true or false for title search */
type boolgen struct {
	src       rand.Source
	cache     int64
	remaining int
}

func (b *boolgen) Bool() bool {
	if b.remaining == 0 {
		b.cache, b.remaining = b.src.Int63(), 63
	}

	result := b.cache&0x01 == 1
	b.cache >>= 1
	b.remaining--

	return result
}

func New() *boolgen {
	return &boolgen{src: rand.NewSource(time.Now().UnixNano())}
}

/* end of random true false */

//function to check producer for the channel -together with  chaincode instantiate control at network level.
func (c *Mrtgexchg) producer(stub shim.ChaincodeStubInterface) bool {
	creatorByte, _ := stub.GetCreator()
	v, err := stub.GetState(string(creatorByte))
	if err != nil {
		return false
	}
	return string(v) == "producer"
}

// write to different ledgers- records, books and lending
func writeToRecordsLedger(stub shim.ChaincodeStubInterface, re RealEstate, txnType string) pb.Response {

	if txnType != "createRealEstate" {
		//add TransactionHistory, first check if map has been initialized
		_, ok := re.TransactionHistory["createRealEstate"]
		if ok {
			re.TransactionHistory[txnType] = getTimeNow()
		} else {
			return shim.Error("......Records Transaction history is not initialized")
		}
	}
	// Encode JSON data
	reAsBytes, err := json.Marshal(re)

	// Store in the Blockchain
	err = stub.PutState(re.RealEstateID, reAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func writeToLendingLedger(stub shim.ChaincodeStubInterface, mrtg Mortgage, txnType string) pb.Response {

	if txnType != "initiateMortgage" {
		//add TransactionHistory
		//first check if map has been initialized
		_, ok := mrtg.TransactionHistory["initiateMortgage"]
		if ok {
			mrtg.TransactionHistory[txnType] = getTimeNow()
		} else {
			return shim.Error("......Mortgage Transaction history is not initialized")
		}
	}

	fmt.Println("++++++++++++++ writing to lending ledger Mortgage Entry=\n ", txnType, " \n", mrtg)

	// Encode JSON data
	mrtgAsBytes, err := json.Marshal(mrtg)

	// Store in the Blockchain
	err = stub.PutState(mrtg.CustID, mrtgAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func writeToBooksLedger(stub shim.ChaincodeStubInterface, bks Books, txnType string) pb.Response {

	if txnType != "initiateBooks" {
		//add TransactionHistory
		//first check if map has been initialized
		_, ok := bks.TransactionHistory["initiateBooks"]
		if ok {
			bks.TransactionHistory[txnType] = getTimeNow()
		} else {
			return shim.Error("......Books Transaction history is not initialized")
		}
	}
	fmt.Println("++++++++++++++ writing to books ledger Books Entry=\n ", txnType, " \n", bks)
	// Encode JSON data
	bksAsBytes, err := json.Marshal(bks)

	// Store in the Blockchain
	err = stub.PutState(bks.RealEstateID, bksAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

/* -------------------------------------------------------------------------------------------------
Implement Init
 -------------------------------------------------------------------------------------------------*/
func (c *Mrtgexchg) Init(stub shim.ChaincodeStubInterface) pb.Response {
	creatorByte, err := stub.GetCreator()
	if err != nil {
		return shim.Error("GetCreator err")
	}
	stub.PutState(string(creatorByte), []byte("producer"))
	return shim.Success(nil)

}

/* -------------------------------------------------------------------------------------------------
 The Invoke method is called as a result of an application request to run the Smart Contract ""
 The calling application program has also specified the particular smart contract function to be called, with arguments
-------------------------------------------------------------------------------------------------*/
func (c *Mrtgexchg) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	function, args := stub.GetFunctionAndParameters() // get function name and args

	if c.producer(stub) { //we are going to restrict creating new blocks on ledgers to those only who instantiate the channel
		if function == "createRealEstate" { //only Registry can create on records ledger
			return c.createRealEstate(stub, args)
		} else if function == "initiateBooks" { //only appraiser can do and write to books ledger
			return c.initiateBooks(stub, args)
		} else if function == "initiateMortgage" { //only bank can do and write to lending ledger
			return c.initiateMortgage(stub, args)
		} else if function == "closeMortgage" { //only bank can do and write to lending ledger
			return c.closeMortgage(stub, args)
		} else if function == "recordPurchase" { //only Registry writes to records ledger
			return c.recordPurchase(stub, args)
		}
	}
	if function == "changeTitle" {
		return c.changeTitle(stub, args)
	} else if function == "queryBooks" {
		return c.queryBooks(stub, args)
	} else if function == "queryLending" {
		return c.queryLending(stub, args)
	} else if function == "getTitle" {
		return c.getTitle(stub, args)
	} else if function == "getFicoScores" {
		return c.getFicoScores(stub, args)
	} else if function == "getAppraisal" {
		return c.getAppraisal(stub, args)
	} else if function == "getInsuranceQuote" {
		return c.getInsuranceQuote(stub, args)
	} else if function == "query" {
		return c.query(stub, args)
	} else if function == "queryAll" {
		return c.queryAll(stub)
	}

	return shim.Error("+~+~+~+~+No matching chain code function found-- create, initiate, close and record mortgage can only be invoked by chaincode instantiators which are Bank, Registry and Appraiser+~+~+~+~+~+~+~+~")

}

/* -------------------------------------------------------------------------------------------------
createRealEstate puts  real estate in the records Blockchain
RealEstateID         string // This one will be our key
Address              string
Value                float64
Details              string
Owner                string
TransactionHistory map[string]string
 -------------------------------------------------------------------------------------------------*/
func (c *Mrtgexchg) createRealEstate(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 5 {
		return shim.Error("createRealEstate arguments usage: RealEstateID, Address, Value, Details, Owner")
	}

	TransactionHistory := make(map[string]string)
	TransactionHistory["createRealEstate"] = getTimeNow()

	// A newly created property is available
	value, _ := strconv.ParseFloat(args[2], 64)
	re := RealEstate{args[0], args[1], value, args[3], args[4], TransactionHistory}

	writeToRecordsLedger(stub, re, "createRealEstate")
	return shim.Success(nil)
}

/* -------------------------------------------------------------------------------------------------
recordPurchase puts  real estate in the records Blockchain with updated owner
RealEstateID         string // This one will be our key
Address              string
Value                float64
Details              string
Owner                string
TransactionHistory map[string]string
-------------------------------------------------------------------------------------------------*/
func (c *Mrtgexchg) recordPurchase(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("createRealEstate arguments usage: RealEstateID,  NewOwner")
	}

	// Look for the RealEstateID
	v, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("RealEstateID " + args[0] + " not found ")
	}

	// Get Information from Blockchain
	var re RealEstate
	// Decode JSON data
	json.Unmarshal(v, &re)

	//first we need to invoke chanincode on books channel to get results value New Owner
	callArgs := make([][]byte, 2)

	callArgs[0] = []byte(QueryBooksString)
	callArgs[1] = []byte(args[0])

	res := stub.InvokeChaincode(BooksChaincode, callArgs, BooksChannel)
	//fmt.Println("************************ received  from books for realestateID=", mrtg.RealEstateID, " Response status=", res.GetStatus(), "payload=", res.Payload)

	var bks Books
	err = json.Unmarshal(res.Payload, &bks)
	if err != nil {
		return shim.Error("Could not Unmarshal aBooks object")
	}

	if bks.NewTitleOwner != "" { //if the new owner is a non blank field then it means the loan was funded and new owner was populated.
		re.Owner = bks.NewTitleOwner
	}
	writeToRecordsLedger(stub, re, "recordPurchase")
	return shim.Success(nil)
}

/* -------------------------------------------------------------------------------------------------
 initiateBooks  initializes books with realestates on records Blockchain
RealEstateID         string // This one will be our key
Appraisal            float64 // this will contain its status on the exchange
NewTitleOwner        string
TransactionHistory map[string]string
-------------------------------------------------------------------------------------------------*/

func (c *Mrtgexchg) initiateBooks(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("createRealEstate arguments usage: RealEstateID")
	}

	TransactionHistory := make(map[string]string)
	TransactionHistory["initiateBooks"] = getTimeNow()

	// A newly created property is available
	books := Books{args[0], 0.0, "", false, TransactionHistory}

	writeToBooksLedger(stub, books, "initiateBooks")
	return shim.Success(nil)
}

/* -------------------------------------------------------------------------------------------------
initiate mortgage
CustID               string // This one will be our key
RealEstateID         string //
LoanAmount           float64
Fico                 float64
Insurance            float64
Appraisal            float64 //this we will get from books ledger
Status               String  //status of the mortgage pending -> closed or not funded due to criteria
TransactionHistory map[string]string
 -------------------------------------------------------------------------------------------------*/

func (c *Mrtgexchg) initiateMortgage(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 3 {
		return shim.Error("initiateMortgage arguments usage: CustID, RealEstateID, LoanAmount")
	}
	TransactionHistory := make(map[string]string)
	TransactionHistory["initiateMortgage"] = getTimeNow()

	amount, _ := strconv.ParseFloat(args[2], 64)
	mrtg := Mortgage{args[0], args[1], amount, 0.0, 0.0, 0.0, "Pending", TransactionHistory}

	writeToLendingLedger(stub, mrtg, "initiateMortgage")
	return shim.Success(nil)
}

/* -------------------------------------------------------------------------------------------------
getFicoScores generates a score and updates the lending ledger
 -------------------------------------------------------------------------------------------------*/
func (c *Mrtgexchg) getFicoScores(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("getFicoScores function needs the customerID as argument")
	}

	v, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("CustomerID " + args[0] + " not found ")
	}

	// Get Information from Blockchain
	var mrtg Mortgage
	// Decode JSON data
	json.Unmarshal(v, &mrtg)

	// update FIco  randomly generated betweenn 600-800
	mrtg.Fico = random(FicoHigh, FicoLow)
	writeToLendingLedger(stub, mrtg, "getFicoScores")
	return shim.Success(nil)
}

/* -------------------------------------------------------------------------------------------------
get appraisal updated by appraiser on books ledger
RealEstateID         float64 // This one will be our key
Appraisal            float64
NewTitleOwner        string
TransactionHistory map[string]string
-------------------------------------------------------------------------------------------------*/
func (c *Mrtgexchg) getAppraisal(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("getAppraisal function needs the RealEstateID as argument")
	}

	// Look for the ID first number
	v, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("RealEstateID " + args[0] + " not found ")
	}

	// Get Information from Blockchain
	var books Books
	// Decode JSON data
	json.Unmarshal(v, &books)

	// update appraisal between 1-2 million
	books.Appraisal = random(AppraisalHigh, AppraisalLow)
	writeToBooksLedger(stub, books, "getAppraisal")
	return shim.Success(nil)
}

/* -------------------------------------------------------------------------------------------------
get appraisal updated by appraiser on books ledger
RealEstateID         float64 // This one will be our key
Appraisal            float64
NewTitleOwner        string
TransactionHistory map[string]string
-------------------------------------------------------------------------------------------------*/
func (c *Mrtgexchg) getTitle(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("getTitle function needs the RealEstateID as argument")
	}

	// Look for the ID first number
	v, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("RealEstateID " + args[0] + " not found ")
	}

	// Get Information from Blockchain
	var books Books
	// Decode JSON data
	json.Unmarshal(v, &books)

	// update Title randonly for true or false
	r := New()
	books.TitleStatus = r.Bool()
	writeToBooksLedger(stub, books, "getTitle")
	return shim.Success(nil)
}

/* -------------------------------------------------------------------------------------------------
get getInsuranceQuote to be called by insurance and updated by insurer on lending ledger
it will need access to customer details, fico and RealEstateID
-------------------------------------------------------------------------------------------------*/
func (c *Mrtgexchg) getInsuranceQuote(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("getInsuranceQuote function needs the customerID and  RealEstateID as argument")
	}

	// Look for the customerID
	v, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("CustomerID " + args[0] + " not found ")
	}

	// Get Information from Blockchain
	var mrtg Mortgage
	// Decode JSON data
	json.Unmarshal(v, &mrtg)

	// update insurance  randomly generated betweenn 2500-5000
	mrtg.Insurance = random(InsuranceHigh, InsuranceLow)
	writeToLendingLedger(stub, mrtg, "getInsuranceQuote")
	return shim.Success(nil)
}

/* -------------------------------------------------------------------------------------------------
changeTitle to be called by bank and updated by appraiser on books ledger
RealEstateID         string // This one will be our key
Appraisal            float64
NewTitleOwner        string
TransactionHistory map[string]string
 -------------------------------------------------------------------------------------------------*/
func (c *Mrtgexchg) changeTitle(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("changeTitle function needs the RealEstateID and new owner as argument")
	}

	// Look for the ID first number
	v, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("RealEstateID " + args[0] + " not found ")
	}

	// Get Information from Blockchain
	var books Books
	// Decode JSON data
	json.Unmarshal(v, &books)

	//first check if the mortgage is funded else reject the title change to new owner
	callArgs := make([][]byte, 2)

	callArgs[0] = []byte(QueryLendingString)
	callArgs[1] = []byte(args[1]) //this is the new customer passed to this function

	res := stub.InvokeChaincode(LendingChaincode, callArgs, LendingChannel)

	var mrtg Mortgage
	err = json.Unmarshal(res.Payload, &mrtg)
	if err != nil {
		return shim.Error("Could not Unmarshal Mortgage object")
	}

	// update owner if mortgage is Funded else it stays blank
	if mrtg.Status == "Funded" {
		books.NewTitleOwner = mrtg.CustID //we will rely on the books leder to update new owner
	}
	writeToBooksLedger(stub, books, "changeTitle")
	return shim.Success(nil)
}

/* -------------------------------------------------------------------------------------------------
closeMortage  updates the lending ledger
-------------------------------------------------------------------------------------------------*/
func (c *Mrtgexchg) closeMortgage(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("closeMortgage function needs the customerID as argument")
	}

	// Look for the serial number
	v, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("CustomerID " + args[0] + " not found ")
	}
	// Get Information from Blockchain
	var mrtg Mortgage
	// Decode JSON data
	err = json.Unmarshal(v, &mrtg)
	if err != nil {
		return shim.Error("Could not Unmarshal results from the Mortgage object")
	}

	//first we need to invoke chanincode on books channel to get appraisal and title search results value of the house provided by Appraiser and Title company
	callArgs := make([][]byte, 2)

	callArgs[0] = []byte(QueryBooksString)
	callArgs[1] = []byte(mrtg.RealEstateID)

	res := stub.InvokeChaincode(BooksChaincode, callArgs, BooksChannel)

	var bks Books
	err = json.Unmarshal(res.Payload, &bks)
	if err != nil {
		return shim.Error("Could not Unmarshal aBooks object")
	}

	// update status of the mortgage to funded if fico score, Insurance and appraisal meets the criteria
	fmt.Println("$^$^$^$^$^$^$^$^$^$^$^$^$^$^$^$^$ Trying to close mortgage loan\n",
		"FicoScore=", mrtg.Fico, "Fico Threshold=", FicoThreshold, "\n",
		"Insurance Quote=", mrtg.Insurance, "Insurance Threshold=", InsuranceThreshold, "\n",
		"Loan Amount=", mrtg.LoanAmount, "Appraised value=", bks.Appraisal, "\n",
		"Title Status=", bks.TitleStatus, "\n",
		"$^$^$^$^$^$^$^$^$^$^$^$^$^$^$^$^$")
	if mrtg.Fico > FicoThreshold && mrtg.Insurance > InsuranceThreshold && bks.Appraisal > mrtg.LoanAmount && bks.TitleStatus == true {
		mrtg.Status = "Funded"
		fmt.Println("@@@@@@@@@@@@@@@@@@ LOan Funded @@@@@@@@@@@@@@@@@@@@@@")
	} else {
		mrtg.Status = "Does not meet criteria for fico and insurance and title an appraised value"
		fmt.Println("--------------------- LOan Rejected------------------------")
	}
	//update lending ledger for appraisal with books appraisal
	mrtg.Appraisal = bks.Appraisal
	writeToLendingLedger(stub, mrtg, "closeMortage")

	return shim.Success(nil)
}

/* -------------------------------------------------------------------------------------------------
// queryRecords, Lending or Books gives all stored keys in the  database- ledger needs to be passed in
 -------------------------------------------------------------------------------------------------*/
func (c *Mrtgexchg) queryAll(stub shim.ChaincodeStubInterface) pb.Response {

	// resultIterator is a StateQueryIteratorInterface
	resultsIterator, err := stub.GetStateByRange("", "")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString("\n,")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}\n")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("\n]")

	fmt.Printf("- queryAll:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/* -------------------------------------------------------------------------------------------------
// queryDetail gives all fields of stored data and needs the key
ledger needs to be passed in
 -------------------------------------------------------------------------------------------------*/
func (s *Mrtgexchg) query(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("query Incorrect number of arguments. Expecting 1")
	}
	asBytes, err := APIstub.GetState(args[0])

	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(asBytes)

}

/* -------------------------------------------------------------------------------------------------
// queryDetail gives all fields of stored data and needs the key
 -------------------------------------------------------------------------------------------------*/
func (s *Mrtgexchg) queryBooks(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("queryBooks Incorrect number of arguments. Expecting 1")
	}
	asBytes, err := APIstub.GetState(args[0])

	if err != nil {
		return shim.Error(err.Error())
	} else {
		var bks Books
		json.Unmarshal(asBytes, &bks)

		writeToBooksLedger(APIstub, bks, QueryBooksString) //log it for audit
		return shim.Success(asBytes)
	}
}

/* -------------------------------------------------------------------------------------------------
// queryDetail gives all fields of stored data and needs the key
 -------------------------------------------------------------------------------------------------*/
func (s *Mrtgexchg) queryLending(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	asBytes, err := APIstub.GetState(args[0])

	if err != nil {
		return shim.Error(err.Error())
	} else {
		var mrtg Mortgage
		json.Unmarshal(asBytes, &mrtg)

		writeToLendingLedger(APIstub, mrtg, QueryLendingString) //log it for audit
		return shim.Success(asBytes)
	}
}

/* -------------------------------------------------------------------------------------------------*/

func main() {
	err := shim.Start(new(Mrtgexchg))
	if err != nil {
		fmt.Printf("Error starting chaincode sample: %s", err)
	}
}
