const samplesRoot = '/samples/'
const createdRoot = 'created/'
const likedRoot = 'liked/'
const path = []

const audio = document.getElementById('audio')
const progress = document.getElementById('progress')
const progressContainer = document.getElementById('progress-container')
const musicContainer = document.getElementById('music-container')
const playBtn = document.getElementById('play')
const title = document.getElementById('title')

google.load("jquery", "1.3.1");
google.setOnLoadCallback(init);

function init() {
    load(samplesRoot, path);
    $('#audio').bind('ended', next);
    $('#next').click(next);
    $('#prev').click(prev);
    audio.addEventListener('timeupdate', updateProgress)
    progressContainer.addEventListener('click', setProgress);

    playBtn.addEventListener('click', () => {
        const isPlaying = musicContainer.classList.contains('play');

        if (isPlaying) {
            pauseSong();
        } else {
            playSong();
        }
    });
}

function load(root, path) {
    const url = root + path.join('/');
    $.ajax({
        url: url,
        dataType: "json",
        success: function (data) {
            listFile(data)
        }
    });
}

function listFile(files) {
    const $b = $('#playlist').empty();

    function addToList(i, f) {
        if (f.Name[0] === '.' || f.Name[0] === ':') return;
        const $d = $('<div></div>').addClass("action-samples-btn");
        f.Path = path.join('/');
        $('<a></a>').text(f.Name).data('file', f)
            .addClass("file")
            .appendTo($d)
            .dblclick(function () {
                let link = document.createElement('a');
                link.setAttribute('href', "./samples/"+f.Name+".mp3");
                link.setAttribute('download', 'download.mp3');
                link.click();
                return false;
            })
            .click(clickFile);
        $('<a></a>').text('\u2764')
            .attr('id', 'like-btn')
            .click(async function () {
                let name = f.Name;
                let formData = new FormData();
                formData.append("sampleName", name)
                await fetch('/likesample/', {
                    method: "POST",
                    body: formData
                });
            })
            .appendTo($d)
        $d.appendTo($b)
    }

    $.each(files, addToList);
}

function likeSample(name) {
    alert()
}

function downloadFile(fileName) {

}

function clickFile(e) {
    const f = $(e.target).data('file');
    const name = f.Name;
    const path = f.Path;
    const url = samplesRoot + path + name;
    $('#playlist a').removeClass('playing');
    $(e.target).addClass('playing');
    loadSong(url, name)
    playSong()
}

function next() {
    const $next = $('#playlist a.playing').next();
    if ($next.length) $next.click();
}

function loadSong(song, songName) {
    title.innerText = songName;
    audio.src = `${song}`;
}

function prev() {
    const $prev = $('#playlist a.playing').prev();
    if ($prev.length) $prev.click();
}

function playSong() {
    musicContainer.classList.add('play');
    playBtn.querySelector('i.fas').classList.remove('fa-play');
    playBtn.querySelector('i.fas').classList.add('fa-pause');

    audio.play();
}

// Pause song
function pauseSong() {
    musicContainer.classList.remove('play');
    playBtn.querySelector('i.fas').classList.add('fa-play');
    playBtn.querySelector('i.fas').classList.remove('fa-pause');

    audio.pause();
}


function updateProgress(e) {
    const {duration, currentTime} = e.srcElement;
    const progressPercent = (currentTime / duration) * 100;
    progress.style.width = `${progressPercent}%`;
}

function setProgress(e) {
    const width = this.clientWidth;
    const clickX = e.offsetX;
    const duration = audio.duration;

    audio.currentTime = (clickX / width) * duration;
}