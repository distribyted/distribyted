GeneralChart.init();
var cacheChart = new SingleBarChart("chart-cache", "Cache disk");

fetchData();
setInterval(function () {
   fetchData();
}, 2000)

function fetchData() {
    fetch('/api/status')
    .then(function (response) {
        if (response.ok) {
            return response.json();
        } else {
            toastError('Error getting data from server. Response: ' + response.status)
        }
    }).then(function (stats) {
        var download = stats.torrentStats.downloadedBytes / stats.torrentStats.timePassed;
        var upload = stats.torrentStats.uploadedBytes / stats.torrentStats.timePassed;

        GeneralChart.update(download, upload);

        cacheChart.update(stats.cacheFilled, stats.cacheCapacity - stats.cacheFilled);
        document.getElementById("down-speed-text").innerText =
            Humanize.ibytes(download, 1024) + "/s";

        document.getElementById("up-speed-text").innerText =
            Humanize.ibytes(upload, 1024) + " /s";
    })
    .catch(function (error) {
        toastError('Error getting status info: ' + error.message)
    });
}