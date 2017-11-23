package main

import (

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"fmt"
	"sync"
	//"runtime"

//	"gitlab.com/defaulteg/api/modules/metal"
	"github.com/defaulteg/api/database"
	"net/http"

	//"gitlab.com/defaulteg/api/utils"


	//"github.com/defaulteg/api/modules"
	"github.com/defaulteg/api/modules/listener"
)

var (
	Db *sql.DB
	err error
)

var waitGroup sync.WaitGroup


func main() {
	//r:=mux.NewRouter()
	//r.HandleFunc("/", HomeHandler).Methods("GET")


	// Init database connection
	if err := database.Init(); err != nil  {
		fmt.Print(err.Error())
		return
	}
	defer database.Instance.Close()


	listener.InvokeApiServiceListener()
	fmt.Println("Listener started...")

	// Start fetcher service
	//modules.StartFetcherService()

	// Wait eternally
	<- make (chan struct{})




}


func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

/*
	runtime.GOMAXPROCS(1)


	waitGroup.Add(2)
	go xx()
	go x()

	waitGroup.Wait()



func xx() {
	defer waitGroup.Done()
	for i:=0;i<60;i++ {
		fmt.Print(i)
		time.Sleep(time.Millisecond * 1000)
	}


}

func x() {
	defer waitGroup.Done()
	for i:=0;i<60;i++ {
		fmt.Print(i)
		time.Sleep(time.Millisecond * 700)
	}

}



*/
