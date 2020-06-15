var fileChunks = new FileChunks();

fetchData();
setInterval(function () {
    fetchData();
}, 2000)

const MB = 1048576;

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
                    var seeders = torrentStat.seeders;
                    var peers = torrentStat.peers;
                    var pieceSize = torrentStat.pieceSize;

                    document.getElementById("up-down-speed-text-" + torrentStat.hash).innerText =
                        Humanize.ibytes(download, 1024) + "/s down, " + Humanize.ibytes(upload, 1024) + "/s up";
                    document.getElementById("peers-seeders-" + torrentStat.hash).innerText =
                        peers + " peers, " + seeders + " seeders."
                    document.getElementById("piece-size-" + torrentStat.hash).innerText = "Piece size: " + Humanize.bytes(pieceSize, 1024)

                    var className = "";

                    if (seeders < 2) {
                        className = "text-danger";
                    } else if (seeders >= 2 && seeders < 4) {
                        className = "text-warning";
                    } else {
                        className = "text-success";
                    }

                    document.getElementById("peers-seeders-" + torrentStat.hash).className = className;

                    if (pieceSize <= MB) {
                        className = "text-success";
                    } else if (pieceSize > MB && pieceSize < (MB * 4)) {
                        className = "text-warning";
                    } else {
                        className =  "text-danger";
                    }

                    document.getElementById("piece-size-" + torrentStat.hash).className = className;
                });
            });
        })
        .catch(function (error) {
            console.log('Error getting status info: ' + error.message);
        });
}