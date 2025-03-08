// main.go
package main

import (
	_ "consumer_api/src/products/application"
	infrastructure "consumer_api/src/products/infraestructure"
	
	_ "github.com/streadway/amqp"
)


func main() {
	//r := gin.Default()
	//r.GET("/web_socket", handleWebSocket)
	infrastructure.Init()
	//r.Run(":3001")
	
}
