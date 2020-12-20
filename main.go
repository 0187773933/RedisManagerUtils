package main

import (
	"fmt"
	"encoding/json"
	manager "github.com/0187773933/RedisManagerUtils/manager"
)

type TestStruct struct {
	wadu string `json:wadu`
	waduagain []int `json:wwaduagain`
}

func main() {
	fmt.Println("main")
	manager := manager.Manager{}
	manager.Connect( "localhost:6379" , 3 , "" )
	//result := manager.Get( "STATE.PREVIOUS.NAME" )
	//result := manager.GetJSON( "testmeta" )
	//fmt.Println( result )
	test_data := TestStruct{
		wadu: "wadu wadu wer234" ,
		waduagain: []int{1,2,3,4} ,
	}
	fmt.Println( test_data )
	json_string , error := json.Marshal( test_data )
	if error != nil { panic( error ) }
	fmt.Println( string( json_string ) )
	// manager.SetJSON( "testmeta" , test_data )
}