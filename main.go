package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
	"github.com/gorilla/websocket"
)

var globalID int64
var dummyQuantities []int
var productsMap = map[string]string{"ftyq4": "Vaporizer", "t3ja2": "Desiccator", "z2lnu": "Humidifier"} //exmple data
var productHashes = []string{"ftyq4", "t3ja2", "z2lnu"}                                                 //example data

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type ClientManager struct {
	clients  map[*websocket.Conn]bool
	orders   []Order
	messages chan string
	lock     sync.Mutex
}

type Order struct {
	ID          int64
	ProductHash string `json:"product_hash"`
	Quantity    int    `json:"quantity"`
	ProductName string
}

func GenerateID() int64 {
	return atomic.AddInt64(&globalID, 1)
}

func getDummyProductData() (string, string) {
	randomIndex := rand.Intn(len(productHashes))
	randomProductHash := productHashes[randomIndex]

	productName, _ := getProductNameFromHash(randomProductHash)
	return randomProductHash, productName
}

func getProductNameFromHash(productHash string) (string, error) {
	productName, exists := productsMap[productHash]
	if !exists {
		return "", fmt.Errorf("product hash %s not found", productHash)
	}

	return productName, nil
}

func getDummyQuantity() int {
	return dummyQuantities[rand.Intn(len(dummyQuantities))]
}

func createClientManager() *ClientManager {
	return &ClientManager{
		clients:  make(map[*websocket.Conn]bool),
		messages: make(chan string, 10),
	}
}

func createOrder() Order {
	return Order{
		ID: GenerateID(),
	}
}

func createDummyOrder() Order {
	order := createOrder()
	order.ProductHash, order.ProductName = getDummyProductData()
	order.Quantity = getDummyQuantity()

	if order.ProductHash == "" || order.ProductName == "" || order.Quantity <= 0 {
		logError("Invalid dummy order", fmt.Errorf("invalid data: %v", order))
	}

	return order
}

func (cm *ClientManager) createMessageFromOrder(order Order) {
	message := fmt.Sprintf("There is a new order for: %s in the quantiy of %d", order.ProductName, order.Quantity)

	cm.lock.Lock()
	cm.messages <- message
	cm.lock.Unlock()
}

func (cm *ClientManager) handleConnection(w http.ResponseWriter, r *http.Request) {

	connectedClient, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		logError("Error upgrading to WebSocket", err)
		return
	}

	defer connectedClient.Close()
	fmt.Println("Client Connected!")

	cm.lock.Lock()
	cm.clients[connectedClient] = true
	cm.lock.Unlock()

	err = connectedClient.WriteJSON(productsMap)
	if err != nil {
		logError("Error writing JSON with products to client", err)
		return
	}

	for {
		_, msg, err := connectedClient.ReadMessage()
		if err != nil {
			logError("Error reading message from client", err)
			break
		}

		cm.processUserMessage(msg)
	}
}

func logError(context string, err error) {
	if err != nil {
		log.Printf("[%s] ERROR: %s - %s\n", time.Now().Format(time.RFC3339), context, err)
	}
}

func (cm *ClientManager) processUserMessage(msg []byte) {
	order := createOrder()
	err := json.Unmarshal(msg, &order)
	if err != nil {
		logError("JSON decoding error", err)
		return
	}

	if order.ProductHash == "" {
		logError("Missing product hash", fmt.Errorf("product hash is empty"))
		return
	}

	order.ProductName, _ = getProductNameFromHash(order.ProductHash)

	cm.processOrder(order)
}

func (cm *ClientManager) processOrder(order Order) {
	defer cm.createMessageFromOrder(order)
	cm.orders = append(cm.orders, order)
}

func (cm *ClientManager) sendMessages() {
	for message := range cm.messages {
		cm.lock.Lock()
		for client := range cm.clients {
			err := client.WriteJSON(message)

			if err != nil {
				logError("Error writing message to client", err)
				cm.removeClient(client)
			}
		}
		cm.lock.Unlock()
	}
}

func (cm *ClientManager) removeClient(client *websocket.Conn) {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	err := client.Close()
	if err != nil {
		logError("Error closing WebSocket connection", err)
	}

	delete(cm.clients, client)

}

func (cm *ClientManager) addDummyOrders() { //simulates regular order placement via dummies
	for {
		time.Sleep((5 * time.Second))
		cm.processOrder(createDummyOrder())
	}
}

func init() {
	quantities := make([]int, 10)
	for i := 0; i < 10; i++ {
		quantities[i] = i + 1
	}

	dummyQuantities = quantities
}

func main() {
	orderManager := createClientManager()
	http.HandleFunc("/ws", orderManager.handleConnection)

	fmt.Println("WebSocket server started on ws://localhost:8080/ws")

	go orderManager.addDummyOrders()
	go orderManager.sendMessages()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		logError("Server failed to start", err)
	}
}
