window.addEventListener("load", function (evt) {
    const socket = new WebSocket('ws://localhost:9000/chat');
    
    var url = window.location.pathname
    var username = url.split("chat/")[1];

    socket.addEventListener('open', function() {
        msg = JSON.stringify(Message(username, "handshake from client", "HANDSHAKE"))
        socket.send(msg)
        addUser(username)
    })
    
    socket.addEventListener('message', function(event) {
        addMessage(event.data)
    })

    socket.addEventListener('close', function (event) {
        socket.close()
        removeUser(username)
    })
    
    function Message(username, message, message_type) {
        return {
            username: username,
            content: message,
            message_type: message_type
        }
    }
    
    function addUser(username) {
        var userListParent = document.getElementById("chatbox__user-list")
    
        var new_user = document.createElement('div');
        new_user.innerHTML = `
        <div class='chatbox__user--active' id="`+ username +`">
            <p>`+ username +`</p>
        </div>`
    
        userListParent.appendChild(new_user)
    }
    
    function removeUser(username) {
        document.getElementById(username).remove()
    }
    
    function addMessage(msg) {
        message = JSON.parse(msg)
        console.log(message);
    
        var msgBox = document.getElementById("chatbox__messages")
    
        var new_message = document.createElement('div');
        new_message.innerHTML =`
        <div class="chatbox__messages__user-message" id="`+ message.id +`">
            <div class="chatbox__messages__user-message--ind-message">
                <p class="name">`+ message.username +`</p>
                <br/>
                <p class="message">` + message.content + `</p>
            </div>
        </div>
        `
        msgBox.appendChild(new_message)
    }
    
    function removeMessage(id) {
        document.getElementById(id).remove()
    }
})
