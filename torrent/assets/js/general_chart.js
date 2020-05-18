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
        var ctx = document.getElementById('chart-general-network').getContext('2d');
        this._chart = new Chart(ctx, {
            type: 'line',
            data: {
                datasets: [
                    {
                        label: 'Download Speed',
                        fill: false,
                        backgroundColor: 'rgb(255, 99, 132)',
                        borderColor: 'rgb(255, 99, 132)',
                        borderWidth: 1,
                        data: this._downloadData,

                    },
                    {
                        label: 'Upload Speed',
                        fill: false,
                        borderWidth: 1,
                        data: this._uploadData,

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
    }
}