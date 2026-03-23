GeneralChart.init();

Distribyted.dashboard = {
    _cacheChart: new CacheChart("main-cache-chart", "Cache disk"),
    loadView: function () {
        fetch('/api/status')
            .then(function (response) {
                if (response.ok) {
                    return response.json();
                } else {
                    Distribyted.message.error('Error getting data from server. Response: ' + response.status)
                }
            }).then(function (stats) {
                var download = stats.torrentStats.downloadedBytes / stats.torrentStats.timePassed;
                var upload = stats.torrentStats.uploadedBytes / stats.torrentStats.timePassed;

                GeneralChart.update(download, upload);

                Distribyted.dashboard._cacheChart.update(stats.cacheFilled, stats.cacheCapacity - stats.cacheFilled);

                document.getElementById("general-download-speed").innerText =
                    Humanize.ibytes(download, 1024) + "/s";

                document.getElementById("general-upload-speed").innerText =
                    Humanize.ibytes(upload, 1024) + "/s";

                document.getElementById("stat-total-torrents").innerText = stats.torrentStats.totalTorrents;
                document.getElementById("stat-total-routes").innerText = stats.torrentStats.totalRoutes;
                document.getElementById("stat-total-downloaded").innerText = Humanize.bytes(stats.torrentStats.totalDownloadedBytes, 1024);
                document.getElementById("stat-total-uploaded").innerText = Humanize.bytes(stats.torrentStats.totalUploadedBytes, 1024);
            })
            .catch(function (error) {
                Distribyted.message.error('Error getting status info: ' + error.message)
            });
    }
}
