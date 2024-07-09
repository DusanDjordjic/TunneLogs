console.log("test");

const socket = new WebSocket('ws://localhost:8080/connect/test/client');

// Event handler for when the WebSocket connection is established
socket.onopen = function(event) {
    console.log('WebSocket connected.');
};

// Event handler for incoming messages from the server
socket.onmessage = function(event) {
    console.log('Received message from server:', event.data);
};

// Event handler for WebSocket errors
socket.onerror = function(error) {
    console.error('WebSocket error:', error);
};

// Event handler for WebSocket connection closure
socket.onclose = function(event) {
    if (event.wasClean) {
        console.log(`WebSocket connection closed cleanly, code=${event.code} reason=${event.reason}`);
    } else {
        console.error('WebSocket connection died');
    }
};
