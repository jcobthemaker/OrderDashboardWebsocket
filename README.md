# WebSocket Order Management Server

This project implements a WebSocket-based order management system, allowing clients to connect to the server and receive updates about orders placed for different products. The server handles client connections, processes orders, and sends updates to connected clients in real-time. It was aso my chance to right proper README so plaease enjoy.

## Features
- **WebSocket Connections**: Clients can connect to the server via WebSocket and receive real-time product order updates.
- **Dummy Order Generation**: Periodically generates dummy orders to simulate real-world usage.
- **Real-time Updates**: Orders are processed and broadcast to all connected clients in real-time. They are displayed on a very basic frontend just to show PoW.
- **Product Data**: A small set of products is available to choose from in the order process. It is example data with the possibility of adding a connection to the database and as a result saving and reading individual product types

## Technologies Used
- **Go (Golang)**: The backend is built using the Go programming language.
- **WebSocket**: WebSocket protocol is used for real-time communication between clients and the server.
- **Gorilla WebSocket Library**: Used to handle WebSocket connections (`github.com/gorilla/websocket`).
- **JSON**: For serializing and deserializing messages sent between the client and server.
- **JQuery**: Simple HTML with with JQuery in order to display and dinamically handle incoming orders and submitting new ones

## Dependencies
- `github.com/gorilla/websocket` - For managing WebSocket connections and communication.

## How It Works

1. **WebSocket Connection**: 
   Clients connect to the server at the `/ws` endpoint using WebSocket. Upon connection, they receive a list of available products.

2. **Order Processing**: 
   - Clients send orders with product information (product hash and quantity) to the server.
   - The server decodes the order and processes it by adding it to the order queue.
   - The server then broadcasts a message to all connected clients with the details of the new order.

3. **Dummy Order Generation**: 
   The server periodically generates dummy orders (simulating real user orders) and processes them, which are then broadcasted to all clients.

4. **Real-time Updates**: 
   Clients that are connected via WebSocket will receive updates whenever a new order is processed (either a dummy order or a client-submitted order).

## Key Functions

- **`GenerateID()`**:  
  Generates a unique order ID using atomic operations.

- **`createDummyOrder()`**:  
  Simulates creating a dummy order with random data.

- **`handleConnection()`**:  
  Handles WebSocket connections, receives, and processes client messages.

- **`processUserMessage()`**:  
  Processes incoming client messages containing orders.

- **`sendMessages()`**:  
  Broadcasts order updates to all connected clients.

- **`addDummyOrders()`**:  
  Periodically generates and processes dummy orders.

---

## Error Handling

Errors are logged using the **`logError()`** function to capture and report any issues during processing. Common error types include:

- WebSocket connection errors
- Missing or invalid order data
- JSON decoding issues
- Client disconnections

This ensures that the server's behavior is monitored and any issues are recorded with timestamps for troubleshooting.

## Further Development

In the future I plan to add: 
- Custom Error Types
- Connection to cloud or database in order to store product types and user orders
- Adding sync.RWMutex in order to allow multiple clients to read concurrently without deadlocks
- Some function decomposition casue I see that there are still possibilities to break down large functions
- Mechanism of graceful server shutdown
- Worker pool in the scenario of more clients
- Adding some memory optimisations. I was reading about sync.Pool for object reuse but in the time of writing this version I havent worked with this yet. See you in the next update
