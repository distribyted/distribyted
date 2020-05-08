var Humanize = require('./humanize.js');
var Chart = require('chart.js');
var _ = require('chartjs-plugin-stacked100');

var ctx = document.getElementById('chart').getContext('2d');

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

new Chart(document.getElementById("chart-example"), {
    type: "horizontalBar",
    data: {
      labels: ["File"],
      datasets: [
        { label: "bad", data: [25], backgroundColor: "rgba(244, 143, 177, 0.6)" },
        { label: "better", data: [10], backgroundColor: "rgba(255, 235, 59, 0.6)" },
        { label: "good", data: [8], backgroundColor: "rgba(100, 181, 246, 0.6)" },
        { label: "s", data: [25], backgroundColor: "rgba(244, 143, 177, 0.6)" },
      ]
    },
    options: {
      plugins: {
        stacked100: { enable: true }
      }
    }
  });

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