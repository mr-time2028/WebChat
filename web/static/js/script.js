$(document).ready(function(){
    $('#action_menu_btn').click(function(){
        $('.action_menu').toggle();
    });
});

let socket = null
let messageField = document.getElementById("message_input")
let outputField = document.getElementById("message_output")
sendMsgBtn = document.getElementById("send_message")

document.addEventListener("DOMContentLoaded", function() {
    socket = new ReconnectingWebSocket("ws://127.0.0.1:8000/ws", null, {debug: true, reconnectInterval: 3000})
    console.log("connected to the weboscket")

    socket.onmessage = msg => {
        data = JSON.parse(msg.data)
        console.log(data.action)
    }
})
