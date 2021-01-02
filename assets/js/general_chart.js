var GeneralChart = {
    _downloadData: [],
    _uploadData: [],
    _chart: null,
    update: function (download, upload) {
        if (this._downloadData.length > 20) {
            this._uploadData.shift();
            this._downloadData.shift();
        }
        var date = new Date();
        this._downloadData.push({
            x: date,
            y: download,
        });
        this._uploadData.push({
            x: date,
            y: upload,
        });
        this._chart.update();
    },
    init: function () {
        var domElem = document.getElementById('chart-general-network')
        domElem.height = 300;
        var ctx = domElem.getContext('2d');
        this._chart = new Chart(ctx, {
            type: 'line',
            data: {
                datasets: [
                    {
                        label: 'Download Speed',
                        fill: false,
                        backgroundColor: "transparent",
                        borderColor: "rgb(82, 136, 255)",

                        lineTension: 0.3,
                        pointRadius: 5,
                        pointBackgroundColor: "rgba(255,255,255,1)",
                        pointHoverBackgroundColor: "rgba(255,255,255,1)",
                        pointBorderWidth: 2,
                        pointHoverRadius: 8,
                        pointHoverBorderWidth: 1,

                        data: this._downloadData,

                    },
                    {
                        label: 'Upload Speed',
                        fill: false,
                        backgroundColor: "transparent",
                        borderColor: "rgb(82, 136, 180)",

                        lineTension: 0.3,
                        pointRadius: 5,
                        pointBackgroundColor: "rgba(255,255,255,1)",
                        pointHoverBackgroundColor: "rgba(255,255,255,1)",
                        pointBorderWidth: 2,
                        pointHoverRadius: 8,
                        pointHoverBorderWidth: 1,

                        data: this._uploadData,

                    },
                ]
            },
            options: {
                legend: {
                    display: false
                },
                responsive: true,
                maintainAspectRatio: false,
                layout: {
                    padding: {
                        right: 10
                    }
                },
                title: {
                    text: 'Download and Upload speed'
                },
                scales: {
                    xAxes: [{
                        scaleLabel: {
                            display: false,
                        },
                        gridLines: {
                            display: false,
                        },
                        ticks: {
                            display: false,
                        },
                        type: 'time',
                    }],
                    yAxes: [{
                        scaleLabel: {
                            display: false,
                            color: "#eee",
                            zeroLineColor: "#eee",
                        },
                        type: 'linear',
                        ticks: {
                            userCallback: function (tick) {
                                return Humanize.ibytes(tick, 1024) + "/s";
                            },
                            beginAtZero: true
                        },
                    }]
                },
                tooltips: {
                    callbacks: {
                        label: function (tooltipItem, data) {
                            var label = data.datasets[tooltipItem.datasetIndex].label || '';

                            if (label) {
                                label += ': ';
                            }

                            return Humanize.ibytes(tooltipItem.yLabel, 1024) + "/s";
                        }
                    },
                    responsive: true,
                    intersect: false,
                    enabled: true,
                    titleFontColor: "#888",
                    bodyFontColor: "#555",
                    titleFontSize: 12,
                    bodyFontSize: 18,
                    backgroundColor: "rgba(256,256,256,0.95)",
                    xPadding: 20,
                    yPadding: 10,
                    displayColors: false,
                    borderColor: "rgba(220, 220, 220, 0.9)",
                    borderWidth: 2,
                    caretSize: 10,
                    caretPadding: 15
                }
            }
        });
    }
}
