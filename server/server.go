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

var topics = map[string][]brokerClient{}

var payloadCommands = map[string]payloadCommand {
		"PUB": publish,	
}
var nonPayloadCommands = map[string]payloadlessCommand{
	"SUB": subscribe,
}
type payloadCommand func(brokerClient,string, string)


func publish(client brokerClient, topic string, payload string){
	clients,ok := topics[topic]
	if ok {
		for _, element := range clients {
				element.connection.Write([]byte(payload))
		}
	} else {
		fmt.Println("No subscribers!")
	}
	fmt.Println(topic + ": " + payload)
}
type payloadlessCommand func(brokerClient,string)

func subscribe(client brokerClient, topic string){
	value, ok := topics[topic]
	if ok {
		topics[topic] = append(value,client)
	} else {
		slice := make([]brokerClient,1)
		slice[0] = client
		topics[topic] = slice
	}
	fmt.Println("Broker Client: " + client.id + " has subscribed to " + topic)
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
		message := strings.TrimSpace(string(buffer[:messageLength]))
		handleMessage(client,message)
		}
		fmt.Println("A connection is closing!")
		
}

func handleMessage(client brokerClient, message string){
	parts := strings.SplitN(message, " ", 3) // split message into 3 parts, command, topic and payload, where payload is optional.
	
	if len(parts) == 3{
		handlePayloadCommand(client,parts)
	} else if len(parts) == 2{
		handleCommandWithoutPayload(client,parts)
	} else {
		fmt.Println("Unknown Command Type: " + strings.Join(parts," "))
	}

}
func handlePayloadCommand(client brokerClient,parts []string){
	method, ok := payloadCommands[parts[0]]	
	if ok {
		method(client,parts[1],strings.Join(parts[2:]," "))
	}
}
func handleCommandWithoutPayload(client brokerClient,parts []string){
	method,ok := nonPayloadCommands[parts[0]]
	if ok{
		method(client, parts[1])
	}
}
