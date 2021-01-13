function CacheChart(id, name) {
    var ctx = document.getElementById(id).getContext('2d');
    this._chart = new Chart(ctx, {
        type: "doughnut",
        data: {
            labels: ["used", "free"],
            datasets: [
                {
                    label: ["used", "free"],
                    data: [0, 0],
                    backgroundColor: ["#4c84ff", "#8061ef"],
                    borderWidth: 1
                }
            ]
        },
        options: {
            animation: false,
            responsive: true,
            maintainAspectRatio: false,
            legend: {
                display: false
            },
            cutoutPercentage: 75,
            tooltips: {
                titleFontColor: "#888",
                bodyFontColor: "#555",
                titleFontSize: 12,
                bodyFontSize: 14,
                backgroundColor: "rgba(256,256,256,0.95)",
                displayColors: true,
                borderColor: "rgba(220, 220, 220, 0.9)",
                borderWidth: 2
            }
        }
    });

    this.update = function (used, free) {
        this._chart.data.datasets.forEach((dataset) => {
            dataset.data[0] = used;
            if (free < 0) {
                free = 0;
            }
            dataset.data[1] = free;
        });

        this._chart.update();
    };
}