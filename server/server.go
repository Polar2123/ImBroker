package main
import (
        "fmt"
        "net"
        "os"
				"strings"
)

const (


	SERVER_HOST = "localhost"
	SERVER_PORT = "9988"
	SERVER_TYPE = "tcp"	
)

type brokerClient struct {
	id string
	connection net.Conn
}

var payloadCommands = map[string]payloadCommand {
		"PUB": publish,	

	}

type payloadCommand func(string, string)

func publish(topic string, payload string){
	fmt.Println(topic + ": " + payload)
}


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
		var client brokerClient
		client.id = "exampleId"
		client.connection = connection
		go handleBrokerConnection(client)	


}
}
func handleBrokerConnection(client brokerClient){

		fmt.Println("Client Connected")
		buffer := make([]byte, 1024)
		for {
			messageLength, err := client.connection.Read(buffer)
			if err != nil {
				fmt.Println("Had problems reading from buffer!")
				break;	
			}	
		message := string(buffer[:messageLength])
		handleMessage(message)
		}
		fmt.Println("A connection is closing!")
		
}

func handleMessage(message string){
	parts := strings.SplitN(message, " ", 3) // split message into 3 parts, command, topic and payload, where payload is optional.
	
	if len(parts) == 3{
		handlePayloadCommand(parts)
	} else if len(parts) == 2{
		handleCommandWithoutPayload(parts)
	} else {
		fmt.Println("Unknown Command Type: " + strings.Join(parts," "))
	}

}
func handlePayloadCommand(parts []string){
	method, ok := payloadCommands[parts[0]]	
	if ok {
		method(parts[1],parts[2])
	}
}
func handleCommandWithoutPayload(parts []string){}
