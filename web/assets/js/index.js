
window.addEventListener('load', function () {
    document.getElementById("form-username").addEventListener('submit', function (event) {
        var username = document.getElementById("username-input").value
    
        if (username !== "") {
            window.location.href = "/chat/" + username;
        } else {
            alert("username cannot be empty")
        }
    
        event.preventDefault();
    })
})

