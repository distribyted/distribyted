var sizes = ["B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"];

function logn(n, b) {
    return Math.log(n) / Math.log(b);
}

var Humanize = {
    bytes: function (s, base) {
        if (s < 10) {
            return s.toFixed(0) + " B";
        }
        var e = Math.floor(logn(s, base));
        var suffix = sizes[e];
        var val = Math.floor(s / Math.pow(base, e) * 10 + 0.5) / 10;

        var f = val.toFixed(0);

        if (val < 10) {
            f = val.toFixed(1);
        }

        return f + suffix;
    }
};
