package main
import (
        "fmt"
        "net"
        "os"
)

const (

	SERVER_HOST = "localhost"
	SERVER_PORT = "9988"
	SERVER_TYPE = "tcp"

)

func main(){

	fmt.Println("Starting ImBroker server!")
	server, err := net.Listen(SERVER_TYPE,SERVER_HOST+":"+SERVER_PORT);
	if err != nil {
		fmt.Println("Error listening for connections!", err.Error())
		os.Exit(1);
	}
	defer server.Close()

	for {

		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting connection!", err.Error())
		}
		fmt.Println("Client Connected")
		buffer := make([]byte, 1024)

		messageLength, err := connection.Read(buffer)
		if err != nil {

			fmt.Println("Couldn't read message from connection!", err.Error())
		}
		fmt.Println(string(buffer[:messageLength]))
		
	}


}
