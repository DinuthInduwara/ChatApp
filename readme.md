# ChatApp 📨💬

A real-time chat application built with Go and WebSockets, powered by the Gorilla WebSocket library. This project allows multiple users to connect, send messages, and receive real-time updates about user connections and disconnections. 🚀

## Features ✨

-   **Real-time Messaging**: Send and receive messages instantly using WebSockets. 📬
-   **User Management**: Tracks connected users and broadcasts events for new users and disconnections. 👥
-   **Dynamic Port Selection**: Automatically finds a free port if the default (8080) is in use. 🔌
-   **Thread-Safe Operations**: Uses `sync.RWMutex` for safe concurrent access to the user list. 🔒
-   **Simple Web Interface**: Serves static files from the `public` directory for a basic frontend (not included). 🌐

## Prerequisites 📋

Before running the application, ensure you have the following installed:

-   **Go** (version 1.16 or higher) 🐹
-   **Git** (to clone the repository) 📂
-   A modern web browser for testing the WebSocket client (e.g., Chrome, Firefox) 🌍

## Installation ⚙️

1. **Clone the Repository**:

    ```bash
    git clone <repository-url>
    cd ChatApp
    ```

2. **Install Dependencies**: Install the Gorilla WebSocket library:

    ```bash
    go get github.com/gorilla/websocket
    ```

3. **Create a** `public` **Directory** (optional): If you plan to serve a frontend, create a `public` directory in the project root and add your static files (e.g., `index.html`, `script.js`).

    ```bash
    mkdir public
    ```

4. **Run the Application**: Compile and run both `main.go` and `room.go`:

    ```bash
    go run main.go room.go
    ```

    The server will start on port `8080` or a random free port if `8080` is in use. Check the console for the port number (e.g., `Starting server on port 8080`). 🚀

## Usage 🚀

1. **Access the WebSocket Endpoint**: Connect to the WebSocket server using a client. The endpoint is:

    ```
    ws://127.0.0.1:<port>/ws?name=<your-username>
    ```

    Replace `<port>` with the port logged by the server and `<your-username>` with a unique name (e.g., `ws://127.0.0.1:8080/ws?name=Alice`).

2. **Test with a WebSocket Client**:

    - Use a tool like `wscat`:

        ```bash
        wscat -c ws://127.0.0.1:8080/ws?name=Alice
        ```

    - Or create a simple HTML/JavaScript frontend in the `public` directory to connect to the WebSocket server.

3. **Send Messages**: Send JSON-formatted messages with the following structure:

    ```json
    {
    	"type": "NEW_MSG_EVENT",
    	"name": "Alice",
    	"payload": {
    		"to": "Bob",
    		"message": "Hello, Bob!"
    	}
    }
    ```

    Supported event types:

    - `NEW_MSG_EVENT`: Send a message to a specific user.
    - `ALL_USERS_EVENT`: Request a list of connected users.

4. **Receive Events**: The server broadcasts events for:

    - New user connections (`NEW_USER_EVENT`).
    - User disconnections (`USER_DISCONNECT_EVENT`).
    - Messages from other users (`NEW_MSG_EVENT`).
    - List of all users (`ALL_USERS_EVENT`).

## Project Structure 📂

```
ChatApp/
├── main.go          # Entry point, sets up HTTP server and WebSocket endpoint 🌐
├── room.go          # Core logic for WebSocket client management and message broadcasting 📨
├── public/          # Directory for static files (e.g., HTML, JS for frontend) 🖼️
└── README.md        # Project documentation (you're here!) 📝
```

## How It Works 🛠️

-   **Server Setup** (`main.go`):

    -   Initializes an HTTP server with the Gorilla WebSocket library.
    -   Handles WebSocket connections at `/ws` with a `name` query parameter.
    -   Dynamically selects a free port if `8080` is in use.
    -   Serves static files from the `public` directory.

-   **WebSocket Logic** (`room.go`):

    -   Manages connected clients in a thread-safe `Users` map.
    -   Processes incoming messages and broadcasts events (new user, disconnection, messages).
    -   Supports direct messaging between users and user list retrieval.

-   **Port Management**:

    -   The `checkPort` function checks if a port is free and falls back to a random port (1024–65535) if needed.
    -   Updates the `checkOrigin` function to allow connections from the selected port.

## Example Client (JavaScript) 🌐

To test the WebSocket server, create a `public/index.html` file:

```html
<!DOCTYPE html>
<html>
	<head>
		<title>ChatApp</title>
	</head>
	<body>
		<h1>ChatApp</h1>
		<input id="username" placeholder="Enter username" />
		<button onclick="connect()">Connect</button>
		<div id="messages"></div>
		<script>
			let ws;
			function connect() {
				const name = document.getElementById("username").value;
				ws = new WebSocket(`ws://127.0.0.1:8080/ws?name=${name}`);
				ws.onmessage = (event) => {
					document.getElementById(
						"messages"
					).innerHTML += `<p>${event.data}</p>`;
				};
				ws.onopen = () => console.log("Connected!");
				ws.onclose = () => console.log("Disconnected!");
			}
		</script>
	</body>
</html>
```

Access it at `http://127.0.0.1:<port>/fs/index.html`.

## Troubleshooting 🐞

-   **Port Conflict**: If you see `listen tcp :8080: bind: Only one usage of each socket address`, the `checkPort` function will automatically select a new port. Check the console for the port number.
-   **WebSocket Connection Fails**:
    -   Ensure the `name` query parameter is provided.
    -   Verify the port in the WebSocket URL matches the server’s port.
    -   Update the `checkOrigin` function if testing from a different host/port.
-   **Dependencies**: Run `go mod init chatapp` and `go mod tidy` if you encounter dependency issues.

## Contributing 🤝

Want to improve ChatApp? Contributions are welcome! 🎉

1. Fork the repository.
2. Create a feature branch (`git checkout -b feature/awesome-feature`).
3. Commit your changes (`git commit -m "Add awesome feature"`).
4. Push to the branch (`git push origin feature/awesome-feature`).
5. Open a Pull Request.

## License 📜

This project is licensed under the MIT License. See the LICENSE file for details.

## Contact 📧

Have questions? Reach out via GitHub Issues or create a pull request! 😊
