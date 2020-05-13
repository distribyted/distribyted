import Humanize from './humanize.js';
import Chart from 'chart.js';
import _ from 'chartjs-plugin-stacked100';

var ctx = document.getElementById('chart-general-network').getContext('2d');

var downloadData = [];
var uploadData = [];
var labels = [];

var chart = new Chart(ctx, {
    type: 'line',
    labels: labels,
    data: {
        datasets: [
            {
                label: 'Download Speed',
                fill: false,
                backgroundColor: 'rgb(255, 99, 132)',
                borderColor: 'rgb(255, 99, 132)',
                borderWidth: 1,
                data: downloadData,

            },
            {
                label: 'Upload Speed',
                fill: false,
                borderWidth: 1,
                data: uploadData,

            },
        ]
    },
    options: {
        title: {
            text: 'Download and Upload speed'
        },
        scales: {
            xAxes: [{
                scaleLabel: {
                    display: true,
                    labelString: 'Date'
                },
                type: 'time',
                time: {
                    //	parser: timeFormat,
                    tooltipFormat: 'll HH:mm'
                },

            }],
            yAxes: [{
                scaleLabel: {
                    display: true,
                    labelString: 'value'
                },
                type: 'linear',
                ticks: {
                    userCallback: function (tick) {
                        return Humanize.bytes(tick, 1024) + "/s";
                    },
                    beginAtZero: true
                },
            }]
        },
    }
});

var piecesData = []
var fileChart = new Chart(document.getElementById("chart-file-chunks"), {
    type: "horizontalBar",
    data: {
        labels: ["File"],
        datasets: piecesData,
    },
    options: {
        animation: {
            duration: 0
        },
        plugins: {
            stacked100: { enable: true }
        }
    }
});


setInterval(function () {
    fetch('/api/status/852299c530aaed8fa06bdf32d9bd909e0bb76fe7')
        .then(function (response) {
            if (response.ok) {
                return response.json();
            } else {
                console.log('Error getting data from server. Response: ' + response.status);
            }
        }).then(function (stats) {
            piecesData.length = 0;
            stats.PieceChunks.forEach(element => {
                var label, color;
                switch (element.Status) {
                    case "H":
                        label = "checking";
                        color = "#8a5999";
                        break;
                    case "P":
                        label = "partial";
                        color = "#be9600";
                        break;
                    case "C":
                        label = "complete";
                        color = "#208f09";
                        break;
                    case "W":
                        label = "waiting";
                        color = "#8a5999";
                        break;
                    case "?":
                        label = "error";
                        color = "#ff5f5c";
                        break;
                    default:
                        label = "unknown";
                        color = "gray";
                        break;
                }
                piecesData.push({
                    label: label,
                    data: [element.NumPieces],
                    backgroundColor: color,
                });

            });
            fileChart.update();
        })
}, 2000)

setInterval(function () {
    fetch('/api/status')
        .then(function (response) {
            if (response.ok) {
                return response.json();
            } else {
                console.log('Error getting data from server. Response: ' + response.status);
            }
        }).then(function (stats) {
            if (downloadData.length > 20) {
                uploadData.shift();
                downloadData.shift();
                labels.shift();
            }

            var date = new Date();
            downloadData.push({
                x: date,
                y: stats.torrentStats.DownloadedBytes / stats.torrentStats.TimePassed,
            });
            uploadData.push({
                x: date,
                y: stats.torrentStats.UploadedBytes / stats.torrentStats.TimePassed,
            });
            labels.push(date);
            chart.update();
        })
        .catch(function (error) {
            console.log('Error getting status info: ' + error.message);
        });
}, 2000)