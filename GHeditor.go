package main

import (
  "fmt"
  "encoding/json"
  "strings"
  "database/sql"
  "bufio"
  "os"
  "os/user"
  "strconv"
  "runtime"
  _"github.com/mattn/go-sqlite3"
)

func input(reader *bufio.Reader, message string)string{
  fmt.Printf(message)
  line, err := reader.ReadString('\n')
  if err != nil{
    fmt.Println("Reader not possible!")
  }
  line = strings.Replace(line, "\r", "", -1)
  line = strings.Replace(line, "\n", "", -1)
  return line
}

func dbOpen(db *sql.DB)map[string]string{
  rows, err := db.Query("SELECT Transactions, User FROM BankAccounts")
  if err != nil{
    fmt.Println("Query failed!")
  }
  var dbData = map[string]string{}
  var index int
  var transactions string
  var user string
  for rows.Next(){
    rows.Scan(&transactions, &user)
    dbData[user] = transactions
    index++
  }
  return dbData
}

func dbWrite(value string, user string, db *sql.DB){
  db.Exec("UPDATE BankAccounts SET Transactions=? WHERE User=? ", value, user)
}

func getProgramPath()string{
  path, err := os.Executable()
  if err != nil{
    fmt.Println("Executable not found please report this issue.")
    os.Exit(1)
  }

  var programPath string
  var seperator string = "/"

  if runtime.GOOS == "windows"{
    seperator = `\`
  }

  pathList := strings.Split(path, seperator)
  pathList = pathList[1:len(pathList)-1]

  for _, i := range pathList{
    programPath += seperator
    programPath += i
  }
  return programPath+seperator
}

func main(){
  programPath := getProgramPath()
  reader := bufio.NewReader(os.Stdin)

  type dbStructure struct{
    Account string  `json:"account"`
    Money float64  `json:"dinero"`
    Transactions []struct{
      Sender string `json:"cuenta"`
      Amount float64 `json:"cantidad"`
      Reason string `json:"motivo"`
      Date string `json:"fecha"`
    } `json:"transacciones"`
  }
  var dbStruct dbStructure

  User, err := user.Current()
  UserName := User.Username
  files := []string{programPath+"GreyHackDB.db","/home/"+UserName+"/.steam/steam/steamapps/common/Grey Hack/Grey Hack_Data/GreyHackDB.db", "C:/Steam/steamapps/common/Grey Hack/Grey Hack_Data/GreyHackDB.db", "D:/Steam/steamapps/common/Grey Hack/Grey Hack_Data/GreyHackDB.db"}
  var file string

  for _, i := range files{
    fileObj, err := os.Open(i)
    if os.IsNotExist(err) {
      continue
    } else {
      file = i
      fileObj.Close()
      break
    }
  }
  if file == ""{
    fmt.Println("Greyhack database not found, put the GreyHack Editor in the same directory as the database.")
    os.Exit(1)
  }
  fmt.Printf("Using database: %s\n\n", file)

  db, err := sql.Open("sqlite3", file)
  dbData := dbOpen(db)

  for user, _ := range dbData{
    fmt.Println("Account: "+user)
  }

  input_account := input(reader, "\nSelect account: ")

  if dbData[input_account] != ""{
    fmt.Println("Bank account found!\n")

    data := dbData[input_account]
    err = json.Unmarshal([]byte(data), &dbStruct)
    if err != nil{
      fmt.Printf("Error: could not read %s\n", dbData[input_account])
      os.Exit(1)
    }

    fmt.Printf("Account %s has $%v\n", dbStruct.Account, dbStruct.Money)

    input_amount, _ := strconv.ParseFloat(input(reader, "Amount you want to add: "), 64)
    dbStruct.Money += input_amount

    out, _ := json.Marshal(dbStruct)
    dbWrite(string(out), dbStruct.Account, db)

    dbData = dbOpen(db)
    data = dbData[input_account]
    err = json.Unmarshal([]byte(data), &dbStruct)

    fmt.Printf("\nAccount %s has now $%v\n", dbStruct.Account, dbStruct.Money)
    db.Close()
  } else {fmt.Println("Bank account not found!")}
}
