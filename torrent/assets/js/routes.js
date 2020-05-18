var fileChunks = new FileChunks();

fetchData();
setInterval(function () {
    fetchData();
}, 2000)

function fetchData() {
    fetch('/api/routes')
        .then(function (response) {
            if (response.ok) {
                return response.json();
            } else {
                console.log('Error getting data from server. Response: ' + response.status);
            }
        }).then(function (routes) {
            routes.forEach(route => {
                route.torrentStats.forEach(torrentStat => {
                    fileChunks.update(torrentStat.pieceChunks, torrentStat.totalPieces, torrentStat.hash);

                    var download = torrentStat.downloadedBytes / torrentStat.timePassed;
                    var upload = torrentStat.uploadedBytes / torrentStat.timePassed;
                    document.getElementById("up-down-speed-text-" + torrentStat.hash).innerText =
                        Humanize.bytes(download, 1024) + "/s down, " + Humanize.bytes(upload, 1024) + "/s up";

                });
            });
        })
        .catch(function (error) {
            console.log('Error getting status info: ' + error.message);
        });
}