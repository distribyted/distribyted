function FileChunks() {
    this.update = function (chunks, totalPieces, hash) {
        var dom = document.getElementById("file-chunks-" + hash);
        dom.innerHTML = "";
        chunks.forEach(chunk => {
            dom.appendChild(getPrintedChunk(chunk.status, chunk.numPieces, totalPieces));
        });
    };
    var pieceStatus = {
        "H": { class: "bg-warning", tooltip: "checking pieces" },
        "P": { class: "bg-info", tooltip: "" },
        "C": { class: "bg-success", tooltip: "downloaded pieces" },
        "W": { class: "bg-transparent" },
        "?": { class: "bg-danger", tooltip: "erroed pieces" },
    };

    var getPrintedChunk = function (status, pieces, totalPieces) {
        var percentage = totalPieces * pieces / 100;
        var pcMeta = pieceStatus[status]
        var pieceStatusClass = pcMeta.class;
        var pieceStatusTip = pcMeta.tooltip;

        var div = document.createElement("div");
        div.className = "progress-bar " + pieceStatusClass;
        div.setAttribute("role", "progressbar");

        if (pieceStatusTip) {
            div.setAttribute("data-toggle", "tooltip");
            div.setAttribute("data-placement", "top");
            div.setAttribute("title", pieceStatusTip);
        }

        div.style.cssText = "width: " + percentage + "%";

        return div;
    };
};


