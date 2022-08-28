const input = document.createElement('choose-file');
input.type = 'file';

input.onchange = ev => {
    var file = ev.target.files[0];

    // setting up the reader
    var reader = new FileReader();
    reader.readAsText(file,'UTF-8');

    // here we tell the reader what to do when it done reading...
    reader.onload = readerEvent => {
        var content = readerEvent.target.result; // this is the content!
        console.log( content );
    }
}
