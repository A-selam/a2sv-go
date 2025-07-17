package main

import (
	"task_manager/data"
	"task_manager/router"
)

func main(){
	err := data.InitMongo()
	if err != nil{
		return 
	}
	defer data.CloseMongo()
	r := router.InitRouter()
	r.Run("localhost:3000")
}