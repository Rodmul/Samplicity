const file = document.getElementById('choose-file');

async function uploadFile() {
    let name = document.getElementById('name').value;
    let pitch = document.getElementById('pitch').value;
    let start = document.getElementById('start').value;
    let end = document.getElementById('end').value;
    let formData = new FormData();
    formData.append("myFile", file.files[0]);
    formData.append("fileName", name);
    formData.append("pitch", pitch);
    formData.append("start", start);
    formData.append("end", end);
    await fetch('/upload/', {
        method: "POST",
        body: formData
    });
}
