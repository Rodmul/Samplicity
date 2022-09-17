async function uploadData() {
    let username = document.getElementById('username').value;
    let email = document.getElementById('email_input').value;
    let password = document.getElementById('password').value;
    let formData = new FormData();
    formData.append("username", username);
    formData.append("email", email);
    formData.append("password", password);
    await fetch('/login/registration/', {
        method: "POST",
        body: formData
    });
}