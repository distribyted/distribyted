GeneralChart.init();
var cacheChart = new SingleBarChart("chart-cache", "Cache disk");
var workerChart = new SingleBarChart("chart-workers", "Workers");

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
            console.log('Error getting data from server. Response: ' + response.status);
        }
    }).then(function (stats) {
        var download = stats.torrentStats.downloadedBytes / stats.torrentStats.timePassed;
        var upload = stats.torrentStats.uploadedBytes / stats.torrentStats.timePassed;

        GeneralChart.update(download, upload);

        cacheChart.update(stats.cacheFilled, stats.cacheCapacity - stats.cacheFilled);
        workerChart.update(0, stats.poolFree - stats.poolCap);
        document.getElementById("down-speed-text").innerText =
            Humanize.bytes(download, 1024) + "/s";

        document.getElementById("up-speed-text").innerText =
            Humanize.bytes(upload, 1024) + " /s";
    })
    .catch(function (error) {
        console.log('Error getting status info: ' + error.message);
    });
}