package main

import (
	"task_manager/router"
)

func main(){
	r := router.InitRouter()
	r.Run("localhost:3000")
}