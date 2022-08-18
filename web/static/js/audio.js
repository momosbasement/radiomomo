'use strict';

// Create an instance
let wavesurfer = {};

// Init & load audio file
document.addEventListener('DOMContentLoaded', function() {
    var wavesurfer = WaveSurfer.create({
        container: '#waveform',
        waveColor: 'violet',
        progressColor: 'purple'
    });

    wavesurfer.on('error', function(e) {
        console.warn(e);
    });

    // Load audio from URL
    wavesurfer.load('http://localhost:8080/api/track/play/31');

    // Play button
    const button = document.querySelector('[data-action="play"]');

    button.addEventListener('click', wavesurfer.playPause.bind(wavesurfer));
});