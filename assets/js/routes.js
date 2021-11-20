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
                return routes;
            })
            .catch(function (error) {
                Distribyted.message.error('Error getting status info: ' + error.message)
            });
    },

    deleteTorrent: function (route, torrentHash) {
        var url = '/api/routes/' + route + '/torrent/' + torrentHash

        return fetch(url, {
            method: 'DELETE'
        })
            .then(function (response) {
                if (response.ok) {
                    Distribyted.message.info('Torrent deleted.')
                    Distribyted.routes.loadView();
                } else {
                    response.json().then(json => {
                        Distribyted.message.error('Error deletting torrent. Response: ' + json.error)
                    })
                }
            })
            .catch(function (error) {
                Distribyted.message.error('Error deletting torrent: ' + error.message)
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

$("#new-magnet").submit(function (event) {
    event.preventDefault();

    let route = $("#route-string :selected").val()
    let magnet = $("#magnet-url").val()

    let url = '/api/routes/' + route + '/torrent'
    let body = JSON.stringify({ magnet: magnet })

    document.getElementById("submit_magnet_loading").style = "display:block"

    fetch(url, {
        method: 'POST',
        body: body
    })
        .then(function (response) {
            if (response.ok) {
                Distribyted.message.info('New magnet added.')
                Distribyted.routes.loadView();
            } else {
                response.json().then(json => {
                    Distribyted.message.error('Error adding new magnet. Response: ' + json.error)
                }).catch(function (error) {
                    Distribyted.message.error('Error adding new magnet: ' + response.status)
                });
            }
        })
        .catch(function (error) {
            Distribyted.message.error('Error adding torrent: ' + error.message)
        }).then(function () {
            document.getElementById("submit_magnet_loading").style = "display:none"
        });
});