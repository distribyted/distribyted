function SingleBarChart(id, name) {
    var ctx = document.getElementById(id).getContext('2d');
    this._used = [];
    this._free = [];
    this._chart = new Chart(ctx, {
        type: 'bar',
        data: {
            labels:[name],
            datasets: [{
                backgroundColor: "gray",
                label: "used",
                data: this._used,
            },
            {
                backgroundColor: "green",
                label: "free",
                data: this._free,
            }],
        },
        options: {
            legend: {
                display: false,
            },
            animation: false,
            scales: {
                xAxes: [{
                    stacked: true
                }],
                yAxes: [{
                    stacked: true,
                    display: true,
                    ticks: {
                        beginAtZero: true,
                    }
                }]
            }
        },
    });

    this.update = function (used, free) {
        this._used.shift();
        this._free.shift();
        this._used.push(used);
        this._free.push(free);

        this._chart.update();
    };
}