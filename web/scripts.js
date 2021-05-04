let gPlaying = false;
let mediaPlayer = undefined;

bar1 = document.getElementById("bar1");
bar2 = document.getElementById("bar2");

bar1.style.display = 'none';
bar2.style.display = 'none';


function onPlayClick() {
    bar1 = document.getElementById("bar1");
    bar2 = document.getElementById("bar2");
    triangle = document.getElementById("triangle");

    if( gPlaying == true ) {
        console.log( "Pause");
        mediaPlayer.pause()
        gPlaying = false;
        bar1.style.display = 'none';
        bar2.style.display = 'none';
        triangle.style.display = 'block';
        
    }
    else {
        console.log( "Play");
        mediaPlayer.play();
        gPlaying = true;
        bar1.style.display = 'block';
        bar2.style.display = 'block';
        triangle.style.display = 'none';
        
    }
   
}

function onNextClick() {
    console.log("Next.");
}

function main() {
    console.log( "Starting...");

    btnPlay = document.getElementById('play');
    btnNext = document.getElementById('next');
    mediaPlayer = document.getElementById('player');

    if( btnPlay != null ) {
        btnPlay.onclick = onPlayClick;
    }

    if( btnNext != null ) {
        btnNext.onclick = onNextClick;
    }

    
}


window.onload = main;