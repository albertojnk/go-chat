window.addEventListener("load", function (evt) {
    const socket = new WebSocket('ws://localhost:9000/chat');
    
    var url = window.location.pathname
    var username = url.split("chat/")[1];

    socket.addEventListener('open', function() {
        var msg = JSON.stringify(Message(username, "handshake from client", "HANDSHAKE"))
        socket.send(msg)
    })
    
    socket.addEventListener('message', function(event) {
        processMessage(event.data)
    })

    socket.addEventListener('close', function (event) {
        socket.close()
    })

    var canPerform = true

    document.getElementById("disconnect-btn").addEventListener('click', function (event) {
        var msg = JSON.stringify(Message(username, "goodbye from client", "GOODBYE"))
        socket.send(msg)
        this.style.display = "none"
        document.getElementById("connect-btn").style.display = "block"
        var input = document.getElementById("msg-input")
        input.disabled = true
        canPerform = false
        input.placeholder = "Disabled until reconnect..."
    })

    document.getElementById("connect-btn").addEventListener('click', function (event) {
        var msg = JSON.stringify(Message(username, "handshake from client", "HANDSHAKE"))
        socket.send(msg)
        window.location.reload()
    })

    
    if (old_messages != "") {
        populateOldMessages(JSON.parse(old_messages)) 
        function populateOldMessages(messages) {
            messages.forEach(msg => {
                if (msg.id != "") {
                    addMessage(msg)
                }
            }); 
        }
    }

    if (connected_users != "") {
        populateConnectedUsers(JSON.parse(connected_users))
        function populateConnectedUsers(users) {
            for (var u in users) {
                if (users[u].username !== username) {
                    if (users[u].username != "") {
                        addUser(users[u].username)

                    }
                }
            }
        }
    }

    document.getElementById("msg-form").addEventListener('submit', function(event) {
        if (canPerform) {
            var msg_content = document.getElementById("msg-input").value
            msg = JSON.stringify(Message(username, msg_content, "MESSAGE"))
            socket.send(msg)
            
            document.getElementById("msg-input").value = ''
            event.preventDefault();
        }
    })
    
    function Message(username, message, message_type) {
        return {
            username: username,
            content: message,
            message_type: message_type
        }
    }
    
    function addUser(username) {
        if (canPerform) {
            var user = document.getElementById(username)
            if (user == null) {
                var userListParent = document.getElementById("chatbox__user-list")
            
                var new_user = document.createElement('div');
                new_user.classList.add('chatbox__user--active')
                new_user.id = username
                new_user.innerHTML = `<p>`+ username +`</p>`
            
                userListParent.append(new_user)
            }
        }
    }

    function addMessage(message) {
        if (canPerform) {

            var msgBox = document.getElementById("chatbox__messages")
            var new_message = document.createElement('div');
            new_message.classList.add("chatbox__messages__user-message")
            new_message.id = message.id
            var img_id = message.id + ':' + message.username
            var img = ''
            if (username == message.username) {
                img = `<img id="` + img_id +`" src="../assets/img/delete-icon.png" alt="delete message"></img>`
                new_message.innerHTML =`
                    <div class="chatbox__messages__user-message--ind-message">
                        <p class="name">`+ message.username +`</p>
                        <br/>
                        <p class="message">` + message.content + `</p>
                        `+ img +`
                    </div>
                    `
                msgBox.appendChild(new_message)

                document.getElementById(img_id).addEventListener('click', function (event) {
                    if (canPerform) {

                        var id_name_array = event.target.id.split(':')
                        // just making sure
                        if (id_name_array[1] == username) {
                            message.message_type = "DELETEMESSAGE"
                            socket.send(JSON.stringify(message))
                        }
                    }
                })
            } else {
                new_message.innerHTML =`
                    <div class="chatbox__messages__user-message--ind-message">
                        <p class="name">`+ message.username +`</p>
                        <br/>
                        <p class="message">` + message.content + `</p>
                        <img style="width:25px;height:25px;visibility:hidden;"></img>
                    </div>
                `
                msgBox.appendChild(new_message)
            }
            
            msgBox.scrollTop = msgBox.scrollHeight;
        }
    }

    function removeUser(username) {
        document.getElementById(username).remove()
    }
    
    function removeMessage(message) {
        if (canPerform) {
            document.getElementById(message.id).remove()
        }
    }

    function processMessage(msg) {
        message = JSON.parse(msg)
        switch (message.message_type) {
            case "GOODBYE":
                removeUser(message.username)
                break;
            case "HANDSHAKE":
                if (message.username != "") {
                    addUser(message.username)
                    populateConnectedUsers(message.clients)
                }
                break;
            case "DELETEMESSAGE":
                removeMessage(message)
                break;
            case "MESSAGE" :
                addMessage(message)
                break;
            default:
                break;
        }
        
    }
})

