package manager

import (
	"fmt"
	"strconv"
	"context"
	redis_lib "github.com/go-redis/redis/v8"
)

// https://pkg.go.dev/github.com/go-redis/redis/v8

type Manager struct {
	Redis *redis_lib.Client
}

func ( manager *Manager ) Connect( address string , db int , password string ) {
	manager.Redis = redis_lib.NewClient( &redis_lib.Options{
		Addr: address ,
		DB: db ,
		Password: password ,
	})
}

func ( manager *Manager ) Get( redis_key string ) ( result string ) {
	result = "failed"
	var ctx = context.Background()
	get_result , get_error := manager.Redis.Get( ctx , redis_key ).Result()
	//if get_error != nil { panic( get_error ) }
	if get_error != nil { fmt.Println( get_error ) } else {
		result = get_result
	}
	return result
}

func ( manager *Manager ) Set( redis_key string , value string ) ( result string ) {
	result = "failed"
	var ctx = context.Background()
	set_result , set_error := manager.Redis.Set( ctx , redis_key , value , 0 ).Result()
	if set_error != nil { fmt.Println( set_error ) } else {
		result = set_result
	}
	return result
}

func ( manager *Manager ) ListIndex( redis_key string , value int64 ) ( result string ) {
	result = "failed"
	var ctx = context.Background()
	index_result , index_result_error := manager.Redis.LIndex( ctx , redis_key , value ).Result()
	if index_result_error != nil { fmt.Println( index_result_error ) } else {
		result = index_result
	}
	return result
}

func ( manager *Manager ) Increment( redis_key string ) ( result string ) {
	result = "failed"
	var ctx = context.Background()
	_ , increment_error := manager.Redis.Incr( ctx , redis_key ).Result()
	if increment_error != nil { fmt.Println( increment_error ) } else {
		result = "success"
	}
	return result
}

func ( manager *Manager ) Decrement( redis_key string ) ( result string ) {
	result = "failed"
	var ctx = context.Background()
	_ , decrement_error := manager.Redis.Decr( ctx , redis_key ).Result()
	if decrement_error != nil { fmt.Println( decrement_error ) } else {
		result = "success"
	}
	return result
}


func ( manager *Manager ) ListPushLeft( redis_key string , value string ) ( result string ) {
	result = "failed"
	var ctx = context.Background()
	_ , l_push_error := manager.Redis.LPush( ctx , redis_key , value ).Result()
	if l_push_error != nil { fmt.Println( l_push_error ) } else {
		result = "success"
	}
	return result
}

func ( manager *Manager ) ListPushRight( redis_key string , value string ) ( result string ) {
	result = "failed"
	var ctx = context.Background()
	_ , r_push_error := manager.Redis.RPush( ctx , redis_key , value ).Result()
	if r_push_error != nil { fmt.Println( r_push_error ) } else {
		result = "success"
	}
	return result
}

func ( manager *Manager ) Publish( redis_key string , value string ) (result string) {
	result = "failed"
	var ctx = context.Background()
	_ , publish_error := manager.Redis.Do( ctx , "PUBLISH" , redis_key , value ).Result()
	if publish_error != nil { fmt.Println( publish_error ) } else {
		result = "success"
	}
	return result
}

func ( manager *Manager ) Subscribe( redis_key string , callback func( msg string ) ) (result string) {
	result = "failed"
	var ctx = context.Background()
	pubsub := manager.Redis.Subscribe( ctx , redis_key )

	// Wait for confirmation that subscription is created before publishing anything.
	_ , err := pubsub.Receive( ctx )
	if err != nil { fmt.Println( err ); return }

	// Channel For Recieving Messages
	ch := pubsub.Channel()
	defer pubsub.Close()
	//var ctx = context.Background()
	_ , err = pubsub.Receive( ctx )
	if err != nil {
		fmt.Println(err)

	}
	// Consume messages.
	for msg := range ch {
		//fmt.Println( msg.Channel , msg.Payload )
		//fmt.Println( msg )
		callback( msg.Payload )
	}
	return
}




func ( manager *Manager ) CircleCurrent( redis_circular_list_key string ) ( result string , index string ) {
	result = "failed"
	index = "0"
	var ctx = context.Background()

	// 1.) Get Length
	circular_list_length , circular_list_length_error := manager.Redis.LLen( ctx , redis_circular_list_key ).Result()
	if circular_list_length_error != nil { panic( circular_list_length_error ) }
	if circular_list_length < 1 { return; }

	// 2.) Get Current Index
	circular_list_key_index_key := fmt.Sprintf( "%s.INDEX" , redis_circular_list_key )
	circular_list_key_index , circular_list_key_index_error := manager.Redis.Get( ctx , circular_list_key_index_key ).Result()
	if circular_list_key_index_error != nil {
		_ , initialize_index_key_error := manager.Redis.Set( ctx , circular_list_key_index_key , "0" , 0 ).Result()
		if initialize_index_key_error != nil { panic( "Could Not Reinitialize Index Variable" ) }
		circular_list_key_index = "0"
	}
	circular_list_index_int , _ := strconv.ParseInt( circular_list_key_index , 0 , 64 )
	circular_list_index_int_64 := int64( circular_list_index_int )

	// 3.) Get Current Value at Index
	current_in_circle , current_in_circle_error := manager.Redis.LIndex( ctx , redis_circular_list_key , circular_list_index_int_64 ).Result()
	if current_in_circle_error != nil { panic( circular_list_key_index_error ) }
	result = current_in_circle
	index = circular_list_key_index
	return
}

func ( manager *Manager ) CirclePrevious( redis_circular_list_key string ) ( result string ) {

	result = "failed"
	var ctx = context.Background()
	//redis := get_redis_connection( "localhost:6379" , 3 , "" )

	// ==================== 1.) Get Length ===================================================================================================================
	circular_list_length_int_64 , circular_list_length_error := manager.Redis.LLen( ctx , redis_circular_list_key ).Result()
	if circular_list_length_error != nil { panic( "Couldn't Get List Length" ) }
	if circular_list_length_int_64 < 1 { panic( "Circular List Length < 1" ) }
	// ==================== 1.) Get Length ===================================================================================================================


	// ==================== 2.) Get Current Index ============================================================================================================
	circular_list_key_index_key := fmt.Sprintf( "%s.INDEX" , redis_circular_list_key )
	circular_list_key_index , circular_list_key_index_error := manager.Redis.Get( ctx , circular_list_key_index_key ).Result()

	// Create .INDEX on the circular list key if it doesnt already exist
	if circular_list_key_index_error != nil {
		circular_list_key_index = "0"
		_ , initialize_index_key_error := manager.Redis.Set( ctx , circular_list_key_index_key , circular_list_key_index , 0 ).Result()
		if initialize_index_key_error != nil { panic( fmt.Sprintf( "Couldn't SET %s" , circular_list_key_index_key  ) ) }
	}
	circular_list_index_int , _ := strconv.ParseInt( circular_list_key_index , 0 , 64 )
	circular_list_index_int_64 := int64( circular_list_index_int )
	// ==================== 2.) Get Current Index ============================================================================================================


	// ==================== 3.) if index == 0 { Set Index To Last Item In List  } else { Decrement Index } ===================================================
	if circular_list_index_int_64 == 0 {
		circular_list_index_int_64 = ( circular_list_length_int_64 - 1 )
		_ , circular_list_set_result_error := manager.Redis.Set( ctx , circular_list_key_index_key , circular_list_index_int_64 , 0 ).Result()
		if circular_list_set_result_error != nil { panic( fmt.Sprintf( "Could not set %s" , circular_list_key_index_key ) ) }
		} else {
			circular_list_index_int_64 = ( circular_list_index_int_64 - 1 )
			_ , circular_list_decr_result_error := manager.Redis.Decr( ctx , circular_list_key_index_key ).Result()
			if circular_list_decr_result_error != nil { panic( fmt.Sprintf( "Couldn't DECR %s" , circular_list_key_index_key  ) ) }
		}
	// ==================== 3.) if index == 0 { Set Index To Last Item In List  } else { Decrement Index } ===================================================

	// ==================== 4.) Get Previous in List @ Updated Index =========================================================================================
	previous_in_circle_list , previous_in_circle_list_error := manager.Redis.LIndex( ctx , redis_circular_list_key , circular_list_index_int_64 ).Result()
	if previous_in_circle_list_error != nil { panic( fmt.Sprintf( "Couldn't LINDEX %s %d" , circular_list_key_index_key , circular_list_index_int_64 ) ) }
	result = previous_in_circle_list
	// ==================== 4.) Get Previous in List @ Updated Index =========================================================================================

	return
}

func ( manager *Manager ) CircleNext( redis_circular_list_key string ) ( result string ) {

	result = "failed"
	var ctx = context.Background()
	//redis := get_redis_connection( "localhost:6379" , 3 , "" )

	// ==================== 1.) Get Length ==================================================================================================================
	circular_list_length_int_64 , circular_list_length_error := manager.Redis.LLen( ctx , redis_circular_list_key ).Result()
	if circular_list_length_error != nil { panic( "Couldn't Get List Length" ) }
	if circular_list_length_int_64 < 1 { panic( "Circular List Length < 1" ) }
	// ==================== 1.) Get Length ==================================================================================================================


	// ==================== 2.) Get Current Index ===========================================================================================================
	circular_list_key_index_key := fmt.Sprintf( "%s.INDEX" , redis_circular_list_key )
	circular_list_key_index , circular_list_key_index_error := manager.Redis.Get( ctx , circular_list_key_index_key ).Result()

	// Create .INDEX on the circular list key if it doesnt already exist
	if circular_list_key_index_error != nil {
		circular_list_key_index = "0"
		_ , initialize_index_key_error := manager.Redis.Set( ctx , circular_list_key_index_key , circular_list_key_index , 0 ).Result()
		if initialize_index_key_error != nil { panic( fmt.Sprintf( "Couldn't SET %s" , circular_list_key_index_key  ) ) }
	}
	circular_list_index_int , _ := strconv.ParseInt( circular_list_key_index , 0 , 64 )
	circular_list_index_int_64 := int64( circular_list_index_int )
	// ==================== 2.) Get Current Index ===========================================================================================================


	// ==================== 3.) if index == ( circular_list_length - 1 ) { Set Index To 0 } else { Increment Index } ========================================
	if circular_list_index_int_64 == ( circular_list_length_int_64 - 1 ) {
		circular_list_index_int_64 = 0
		_ , circular_list_set_result_error := manager.Redis.Set( ctx , circular_list_key_index_key , circular_list_index_int_64 , 0 ).Result()
		if circular_list_set_result_error != nil { panic( fmt.Sprintf( "Could not set %s" , circular_list_key_index_key ) ) }
		} else {
			circular_list_index_int_64 = ( circular_list_index_int_64 + 1 )
			_ , circular_list_decr_result_error := manager.Redis.Incr( ctx , circular_list_key_index_key ).Result()
			if circular_list_decr_result_error != nil { panic( fmt.Sprintf( "Couldn't DECR %s" , circular_list_key_index_key  ) ) }
		}
	// ==================== 3.) if index == ( circular_list_length - 1 ) { Set Index To 0 } else { Increment Index } ========================================

	// ==================== 4.) Get Next in List @ Updated Index ============================================================================================
	next_in_circle_list , next_in_circle_list_error := manager.Redis.LIndex( ctx , redis_circular_list_key , circular_list_index_int_64 ).Result()
	if next_in_circle_list_error != nil { panic( fmt.Sprintf( "Couldn't LINDEX %s %d" , circular_list_key_index_key , circular_list_index_int_64 ) ) }
	result = next_in_circle_list
	// ==================== 4.) Get Next in List @ Updated Index ============================================================================================

	return
}