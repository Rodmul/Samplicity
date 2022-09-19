const file = document.getElementById('choose-file');

async function uploadFile() {
    let name = document.getElementById('name').value;
    let type = document.getElementById('select-type-field').value;
    let start = document.getElementById('start').value;
    let end = document.getElementById('end').value;
    let checkBoxTime = document.querySelector('#time-checkbox');
    let startTime = document.getElementById('start-time').value;
    let endTime = document.getElementById('end-time').value;
    let formData = new FormData();

    if (checkBoxTime.checked) {
        formData.append("start", startTime);
        formData.append("end", endTime);
    } else {
        formData.append("start", start);
        formData.append("end", end);
    }

    formData.append("myFile", file.files[0]);
    formData.append("fileName", name);
    formData.append("type", type);
    await fetch('/upload/', {
        method: "POST",
        body: formData
    });
}
