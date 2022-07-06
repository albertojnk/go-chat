window.addEventListener("load", function (evt) {
    const socket = new WebSocket('ws://localhost:9000/chat');
    
    var url = window.location.pathname
    var username = url.split("chat/")[1];

    socket.addEventListener('open', function() {
        msg = JSON.stringify(Message(username, "handshake from client", "HANDSHAKE"))
        socket.send(msg)
    })
    
    socket.addEventListener('message', function(event) {
        processMessage(event.data)
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
        new_user.classList.add('chatbox__user--active')
        new_user.id = username
        new_user.innerHTML = `<p>`+ username +`</p>`
    
        userListParent.append(new_user)
    }

    function addMessage(message) {
        var msgBox = document.getElementById("chatbox__messages")
        var new_message = document.createElement('div');
        new_message.classList.add("chatbox__messages__user-message")
        new_message.id = message.id
        new_message.innerHTML =`
            <div class="chatbox__messages__user-message--ind-message">
                <p class="name">`+ message.username +`</p>
                <br/>
                <p class="message">` + message.content + `</p>
            </div>
        `
        msgBox.appendChild(new_message)
    }

    function removeUser(username) {
        document.getElementById(username).remove()
    }
    
    function removeMessage(id) {
        var ele = document.getElementById(id)
        ele.remove()
    }

    function processMessage(msg) {
        message = JSON.parse(msg)
        console.log(message);
        switch (message.message_type) {
            case "GOODBYE":
                removeUser(message.username)
                break;
            case "HANDSHAKE":
                addUser(message.username)
                break;
            case "MESSAGEDELETE":
                removeMessage(message.id)
                break;
            case "MESSAGE" :
                addMessage(message)
                break;
            default:
                break;
        }
        
    }
})
