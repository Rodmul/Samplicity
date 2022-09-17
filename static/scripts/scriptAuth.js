async function uploadData() {
    let username = document.getElementById('username').value;
    let password = document.getElementById('password').value;
    let formData = new FormData();
    formData.append("username", username);
    formData.append("password", password);
    await fetch('/auth/authorization/', {
        method: "POST",
        body: formData
    });
}