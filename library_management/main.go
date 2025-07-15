package main

import (
	"library_management/controllers"
	"library_management/services"
)

func main(){
	services.Init()
	controllers.App()
}