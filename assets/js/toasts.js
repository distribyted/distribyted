function toast(msg, color, time) {
  Toastify({
    text: msg,
    duration: time,
    newWindow: true,
    close: true,
    gravity: "top",
    position: "right",
    backgroundColor: color,
    stopOnFocus: true,
  }).showToast();
}

function toastError(msg) {
  toast(msg, "red", 5000);
}

function toastInfo(msg) {
  toast(msg, "green", 5000);
}

function toastMessage(msg) {
  toast(msg, "grey", 10000);
}

var stream = new EventSource("/api/events");
stream.addEventListener("event", function (e) {
  toastMessage(e.data);
});
