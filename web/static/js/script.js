let socket = null
let messageField = document.getElementById("message_input")
let outputField = document.getElementById("message_output")
sendMsgBtn = document.getElementById("send_message")

document.addEventListener("DOMContentLoaded", function() {
    // socket = new ReconnectingWebSocket("ws://127.0.0.1:8000/ws", null, {debug: true, reconnectInterval: 3000})
    socket = new WebSocket("ws://127.0.0.1:8000/ws", ["json"])

    socket.onopen = () => {
        console.log("authenticating");
        // Send the authentication token as the first message
        const token = 'Bearer dafdasgj';
        let jsonData = {
            authorization: token,
            action: "auth",
            username: "MrTime",
        };
        socket.send(JSON.stringify(jsonData));
    };

    socket.onmessage = msg => {
        data = JSON.parse(msg.data)
        console.log(data.message)
        console.log(data.action)
    }
})
