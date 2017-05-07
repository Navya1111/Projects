package AddressBook
import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gocql/gocql"
	"os"
	"encoding/csv"
	"io"
	"strconv"
	"strings"
)

type Address struct{
	Lastname string
	Firstname string
	Email string
	Phone int
}

func createSession(cluster, keyspace string) (*gocql.Session,error){
	session := gocql.NewCluster(cluster)
	session.Keyspace = keyspace
	session.Consistency = gocql.Quorum
	return session.CreateSession()
}

func ImportAddresses(w http.ResponseWriter, req *http.Request){
	var result = ""
	params := mux.Vars(req)
	sourceFile := params["fileName"]
	sourceFile += ".csv"
	mySession,err := createSession("127.0.0.1","example")
	if err!= nil{
		result = err.Error()
	}
	defer mySession.Close()
	file,err :=os.Open(sourceFile)
	if err!= nil{
		result = err.Error()
	}
	reader := csv.NewReader(file)
	var details Address
	for{
		file1,err := reader.Read()
		if err == io.EOF{
			break
		}
		details.Firstname = file1[0]
		details.Lastname = file1[1]
		details.Phone,_ = strconv.Atoi(file1[2])
		details.Email = file1[3]
		err= mySession.Query("insert into addressbook(phone,email,firstname,lastname) VALUES (?,?,?,?)",details.Phone,details.Email,details.Firstname,details.Lastname).Exec()
                if(err != nil){
			result = err.Error()
		}else{
			result = "Addressbook imported Successfully"
		}
	}
	json.NewEncoder(w).Encode(result)
}

func ExportAddressBook(w http.ResponseWriter, req *http.Request){

	var result = ""
	mySession,err := createSession("127.0.0.1","example")
	if err!= nil{
		result = err.Error()
	}else{
		defer mySession.Close()
		myfile,err := os.Create("AddressBook.csv")
		if err!= nil{
			result = err.Error()
		}
		defer myfile.Close()
		writer := csv.NewWriter(myfile)

		iter := mySession.Query("select * from example.addressbook").Iter()
		defer iter.Close()
		row := map[string]interface{}{}
		for iter.MapScan(row){
			var address []string
			address = append(address, row["firstname"].(string))
			address = append(address, row["lastname"].(string))
			address = append(address, strconv.FormatInt(row["phone"].(int64), 10))
			address = append(address, row["email"].(string))
			writer.Write(address)
			row = map[string]interface{}{}
		}
		writer.Flush()
		result = "AddressBook Exported Successfully"
	}
	json.NewEncoder(w).Encode(result)
}

func GetAddress(w http.ResponseWriter, req *http.Request){

	params := mux.Vars(req)
	Phone := params["phone"]
	//Database connect code and Retrieve
	mySession,err := createSession("127.0.0.1","example")
	if err!= nil{
		fmt.Println(err)
	}
	defer mySession.Close()

	var firstname string
	var lastname string
	var phone int
	var email string
	iter := mySession.Query("select * from AddressBook where phone=?",Phone).Iter()
	for {
		row := map[string]interface{}{
			"Firstname": &firstname,
			"Lastname": &lastname,
			"Phone": &phone,
			"Email": &email,
		}
		if !iter.MapScan(row) {
			break
		}
	}
	json.NewEncoder(w).Encode(iter)
}

func CreateAddress(w http.ResponseWriter, req *http.Request){

	var result = ""
	params := mux.Vars(req)
	phone := params["phone"]
	var details Address
	err :=req.ParseForm()
	if(err != nil){
		fmt.Println(err)
	}
	details.Email = req.Form.Get("email")
	details.Lastname = req.Form.Get("lastname")
	details.Firstname = req.Form.Get("firstname")
	details.Phone,_ = strconv.Atoi(phone)
	mySession,err := createSession("127.0.0.1","example")
	if err!= nil{
		result = err.Error()
	}else{
		err= mySession.Query("insert into addressbook(phone,email,firstname,lastname) VALUES (?,?,?,?)",details.Phone,details.Email,details.Firstname,details.Lastname).Exec()

		if(err == nil){
			result = "Address added to AddressBook Successfully"
		} else {
			result = err.Error()
		}
	}
	defer mySession.Close()
	json.NewEncoder(w).Encode(result)
}

func UpdateAddress(w http.ResponseWriter, req *http.Request){

	var query = "UPDATE addressbook SET "
	var result = ""
	params := mux.Vars(req)
	phone := params["phone"]
	err :=req.ParseForm()
	if(err != nil){
		result = err.Error()
	}
	if(req.Form.Get("email") != ""){
		query = query + "email = '" + req.Form.Get("email") + "',"
	}
	if(req.Form.Get("lastname") != ""){
		query = query + "lastname = '" + req.Form.Get("lastname") + "',"
	}
	if(req.Form.Get("firstname") != ""){
		query = query + "firstname = '" + req.Form.Get("firstname") + "',"
	}
	query = strings.TrimRight(query, ",")
	query = query + " WHERE phone = " + phone
	mySession,err := createSession("127.0.0.1","example")
	if err!= nil{
		result = err.Error()
	}else{
		err= mySession.Query(query).Exec()

		if(err == nil){
			result = "Address Updated Successfully"
		} else {
			fmt.Println(err)
			result = err.Error()
		}
	}
	defer mySession.Close()
	json.NewEncoder(w).Encode(result)
}

func DeleteAddress(w http.ResponseWriter, r *http.Request){

	var result = ""
	params := mux.Vars(r)
	Phone := params["phone"]
	mySession,err := createSession("127.0.0.1","example")
	if err!= nil{
		fmt.Println(err)
	}
	defer mySession.Close()
	err = mySession.Query("Delete from addressbook where phone=?",Phone).Exec()
	if err == nil{
		result = "Address Deleted Succesfully"
	} else {
		result = err.Error()
	}
	json.NewEncoder(w).Encode(result)
}


