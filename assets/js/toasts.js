function toastError(msg) {
    Toastify({
        text: msg,
        duration: 5000, 
        newWindow: true,
        close: true,
        gravity: "top", // `top` or `bottom`
        position: 'right', // `left`, `center` or `right`
        backgroundColor: "red",
        stopOnFocus: true, // Prevents dismissing of toast on hover
      }).showToast();
}

function toastInfo(msg) {
    Toastify({
        text: msg,
        duration: 5000, 
        newWindow: true,
        close: true,
        gravity: "top", // `top` or `bottom`
        position: 'right', // `left`, `center` or `right`
        backgroundColor: "green",
        stopOnFocus: true, // Prevents dismissing of toast on hover
      }).showToast();
}