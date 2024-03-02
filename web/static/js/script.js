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
    socket = new ReconnectingWebSocket("ws://127.0.0.1:8000/ws", null, {debug: true, reconnectInterval: 3000})
    // socket = new WebSocket("ws://127.0.0.1:8000/ws", ["json"])

    socket.onopen = () => {
        console.log("authenticating");
        // Send the authentication token as the first message
        let jsonData = {
            authorization: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJsb2NhbGhvc3QiLCJleHAiOjE3MDk0MDA1NzgsImlhdCI6MTcwOTQwMDI3OCwiaXNzIjoibG9jYWxob3N0IiwibmFtZSI6Ik1yVGltZSIsInN1YiI6ImY5MGMxNTE5LTA5YTUtNDMwNi1iMTMzLTAxMDkxODk4MDNmYyIsInR5cCI6IkpXVCJ9.kKkxXi2nd7H9OYyKmkyfVUBIXR4DdbKZ2x5Qohrvn7s",
            action: "auth",
        };
        socket.send(JSON.stringify(jsonData));
    };

    socket.onmessage = msg => {
        data = JSON.parse(msg.data)
        console.log(data.message)
        console.log(data.action)
        console.log(data.connected_users)
    }
})
