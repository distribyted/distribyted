Handlebars.registerHelper("torrent_status", function (chunks, totalPieces) {
    const pieceStatus = {
        "H": { class: "bg-warning", tooltip: "checking pieces" },
        "P": { class: "bg-info", tooltip: "" },
        "C": { class: "bg-success", tooltip: "downloaded pieces" },
        "W": { class: "bg-transparent" },
        "?": { class: "bg-danger", tooltip: "erroed pieces" },
    };
    const chunksAsHTML = chunks.map(chunk => {
        const percentage = totalPieces * chunk.numPieces / 100;
        const pcMeta = pieceStatus[chunk.status]
        const pieceStatusClass = pcMeta.class;
        const pieceStatusTip = pcMeta.tooltip;

        const div = document.createElement("div");
        div.className = "progress-bar " + pieceStatusClass;
        div.setAttribute("role", "progressbar");

        if (pieceStatusTip) {
            div.setAttribute("data-toggle", "tooltip");
            div.setAttribute("data-placement", "top");
            div.setAttribute("title", pieceStatusTip);
        }

        div.style.cssText = "width: " + percentage + "%";

        return div.outerHTML;
    });

    return '<div class="progress mb-3">' + chunksAsHTML.join("\n"); + '</div>'
});

Handlebars.registerHelper("torrent_info", function (peers, seeders, pieceSize) {
    const MB = 1048576;

    var messages = [];

    var errorLevels = [];
    const seedersMsg = "- Number of seeders is too low (" + seeders + ")."
    if (seeders < 2) {
        errorLevels[0] = 2;
        messages.push(seedersMsg);
    } else if (seeders >= 2 && seeders < 4) {
        errorLevels[0] = 1;
        messages.push(seedersMsg);
    } else {
        errorLevels[0] = 0;
    }

    const pieceSizeMsg = "- Piece size is too big (" + Humanize.bytes(pieceSize, 1024) + "). Recommended size is 1MB or less."
    if (pieceSize <= MB) {
        errorLevels[1] = 0;
    } else if (pieceSize > MB && pieceSize < (MB * 4)) {
        errorLevels[1] = 1;
        messages.push(pieceSizeMsg);
    } else {
        errorLevels[2] = 2;
        messages.push(pieceSizeMsg);
    }

    const level = ["text-success", "text-warning", "text-danger"];
    const icon = ["mdi-check", "mdi-alert", "mdi-alert-octagram"];
    const div = document.createElement("div");
    const i = document.createElement("i");

    const errIndex = Math.max(...errorLevels);

    i.className = "mdi " + icon[errIndex];
    i.title = messages.join("\n");

    const text = document.createTextNode(peers + "/" + seeders + " (" + Humanize.bytes(pieceSize, 1024) + " chunks) ");

    div.className = level[errIndex];
    div.appendChild(text);
    div.appendChild(i);

    return div.outerHTML;
});


Distribyted.routes = {
    _template: null,

    _getTemplate: function () {
        if (this._template != null) {
            return this._template
        }

        const tTemplate = fetch('/assets/templates/routes.html')
            .then((response) => {
                if (response.ok) {
                    return response.text();
                } else {
                    Distribyted.message.error('Error getting data from server. Response: ' + response.status);
                }
            })
            .then((t) => {
                return Handlebars.compile(t);
            })
            .catch(error => {
                Distribyted.message.error('Error getting routes template: ' + error.message);
            });

        this._template = tTemplate;
        return tTemplate;
    },

    _getRoutesJson: function () {
        return fetch('/api/routes')
            .then(function (response) {
                if (response.ok) {
                    return response.json();
                } else {
                    Distribyted.message.error('Error getting data from server. Response: ' + response.status)
                }
            }).then(function (routes) {
                // routes.forEach(route => {
                //     route.torrentStats.forEach(torrentStat => {
                //         fileChunks.update(torrentStat.pieceChunks, torrentStat.totalPieces, torrentStat.hash);

                //         var download = torrentStat.downloadedBytes / torrentStat.timePassed;
                //         var upload = torrentStat.uploadedBytes / torrentStat.timePassed;
                //         var seeders = torrentStat.seeders;
                //         var peers = torrentStat.peers;
                //         var pieceSize = torrentStat.pieceSize;

                //         document.getElementById("up-down-speed-text-" + torrentStat.hash).innerText =
                //             Humanize.ibytes(download, 1024) + "/s down, " + Humanize.ibytes(upload, 1024) + "/s up";
                //         document.getElementById("peers-seeders-" + torrentStat.hash).innerText =
                //             peers + " peers, " + seeders + " seeders."
                //         document.getElementById("piece-size-" + torrentStat.hash).innerText = "Piece size: " + Humanize.bytes(pieceSize, 1024)

                //         var className = "";

                //         if (seeders < 2) {
                //             className = "text-danger";
                //         } else if (seeders >= 2 && seeders < 4) {
                //             className = "text-warning";
                //         } else {
                //             className = "text-success";
                //         }

                //         document.getElementById("peers-seeders-" + torrentStat.hash).className = className;

                //         if (pieceSize <= MB) {
                //             className = "text-success";
                //         } else if (pieceSize > MB && pieceSize < (MB * 4)) {
                //             className = "text-warning";
                //         } else {
                //             className =  "text-danger";
                //         }

                //         document.getElementById("piece-size-" + torrentStat.hash).className = className;
                //     });
                // });
                return routes;
            })
            .catch(function (error) {
                Distribyted.message.error('Error getting status info: ' + error.message)
            });
    },

    loadView: function () {
        this._getTemplate()
            .then(t =>
                this._getRoutesJson().then(routes => {
                    document.getElementById('template_target').innerHTML = t(routes);
                })
            );
    }
}