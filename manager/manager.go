package manager

import (
	"context"
	redis_lib "github.com/go-redis/redis/v8"
)

// https://pkg.go.dev/github.com/go-redis/redis/v8

type Manager struct {
	redis *redis_lib.Client
}

func ( manager *Manager ) Connect( address string , db int , password string ) {
	manager.redis = redis_lib.NewClient( &redis_lib.Options{
		Addr: address ,
		DB: db ,
		Password: password ,
	})
}

func ( manager *Manager ) Get( redis_key string ) ( result string ) {
	result = "failed"
	var ctx = context.Background()
	get_result , error := manager.redis.Get( context , redis_key ).Result()
	if error != nil { panic( error ) }
	result = get_result
	return result
}

func ( manager *Manager ) Set( redis_key string , value string ) ( result string ) {
	result = "failed"
	var ctx = context.Background()
	set_result , error := manager.redis.Set( context , redis_key , value , 0 ).Result()
	if error != nil { panic( error ) }
	result = set_result
	return result
}

func ( manager *Manager ) ListIndex( redis_key string , value string ) ( result string ) {
	result = "failed"
	var ctx = context.Background()
	index_result , index_result_error := manager.redis.LIndex( ctx , redis_circular_list_key , circular_list_index_int_64 ).Result()
	if index_result_error != nil { panic( index_result_error ) }
	if error != nil { panic( error ) }
	result = index_result
	return result
}

func ( manager *Manager ) Increment( redis_key string ) ( result string ) {
	result = "failed"
	var context = context.Background()
	_ , increment_error := manager.redis.Incr( ctx , redis_circular_list_key ).Result()
	if increment_error != nil { panic( increment_error ) }
	if error != nil { panic( error ) }
	result = "success"
	return result
}

func ( manager *Manager ) Decrement( redis_key string ) ( result string ) {
	result = "failed"
	var context = context.Background()
	_ , increment_error := manager.redis.Decr( ctx , redis_circular_list_key ).Result()
	if increment_error != nil { panic( increment_error ) }
	if error != nil { panic( error ) }
	result = "success"
	return result
}


func ( manager *Manager ) ListPushLeft( redis_key string ) ( result string ) {
	result = "failed"
	var ctx = context.Background()
	_ , increment_error := manager.redis.LPush( ctx , redis_circular_list_key ).Result()
	if increment_error != nil { panic( increment_error ) }
	if error != nil { panic( error ) }
	result = "success"
	return result
}

func ( manager *Manager ) ListPushRight( redis_key string ) ( result string ) {
	result = "failed"
	var ctx = context.Background()
	_ , increment_error := manager.redis.RPush( ctx , redis_circular_list_key ).Result()
	if increment_error != nil { panic( increment_error ) }
	if error != nil { panic( error ) }
	result = "success"
	return result
}