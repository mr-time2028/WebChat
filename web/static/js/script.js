let socket = null
let messageField = document.getElementById("message_input")
let outputField = document.getElementById("message_output")
sendMsgBtn = document.getElementById("send_message")

window.onbeforeunload = function() {
    console.log("Leaving")
    let jsonData = {};
    jsonData["action"] = "left";
    socket.send(JSON.stringify(jsonData))
}

document.addEventListener("DOMContentLoaded", function() {
    socket = new ReconnectingWebSocket("ws://127.0.0.1:8000/join_room", null, {debug: true, reconnectInterval: 3000})
    // socket = new WebSocket("ws://127.0.0.1:8000/ws", ["json"])

    socket.onopen = () => {
        const payload = {
            token: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJsb2NhbGhvc3QiLCJleHAiOjE3MTA5NTU0MjYsImlhdCI6MTcxMDk1NTEyNiwiaXNzIjoibG9jYWxob3N0IiwibmFtZSI6Ik1yVGltZSIsInN1YiI6IjVjZTJmZjcwLTU1N2EtNDI1ZS04ZTExLWU2YzBlNjk5MzMxMCIsInR5cCI6IkpXVCJ9.TwLADIOkvnbfyWcmSwmy0arHhiN8-8JWySTk8f0SDic",
            room_id: "d864ae89-c6b0-4899-b0f7-4eddce3c53af"
        }

        socket.send(JSON.stringify(payload));
    };

    socket.onmessage = msg => {
        let data = JSON.parse(msg.data)
        console.log(data)
    }
})
