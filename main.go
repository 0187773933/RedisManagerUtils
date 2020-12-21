package main

import (
	"fmt"
	"time"
	"encoding/json"
	redis "github.com/0187773933/RedisManagerUtils/manager"
)

type TestStruct struct {
	Wadu string `json:wadu`
	Waduagain []int `json:waduagain`
}

func main() {

	redis := redis.Manager{}
	redis.Connect( "localhost:6379" , 3 , "" )

	go redis.Subscribe( "LOG.ALL" )

	// Set JSON Example
	test_data := TestStruct{
		Wadu: "wadu wadu wer234" ,
		Waduagain: []int{1,2,3,4} ,
	}
	json_marshal_result , json_marshal_error := json.Marshal( test_data )
	if json_marshal_error != nil { panic( json_marshal_error ) }
	json_string := string( json_marshal_result )
	redis.Set( "testmeta" , json_string )
	fmt.Println( json_string )

	// Get JSON Example
	json_get_test := redis.Get( "testmeta" )
	var json_get_test_struct TestStruct
	json_unmarshal_error := json.Unmarshal( []byte( json_get_test ) , &json_get_test_struct )
	if json_unmarshal_error != nil { panic( json_unmarshal_error ) }
	fmt.Println( json_get_test_struct )

	fmt.Println( redis.CircleCurrent( "test" ) )
	fmt.Println( redis.CircleNext( "test" ) )
	fmt.Println( redis.CircleNext( "test" ) )
	fmt.Println( redis.CircleNext( "test" ) )
	fmt.Println( redis.CirclePrevious( "test" ) )
	fmt.Println( redis.CirclePrevious( "test" ) )
	fmt.Println( redis.CirclePrevious( "test" ) )
	fmt.Println( redis.CircleCurrent( "test" ) )

	time.Sleep( 20*time.Second )

}