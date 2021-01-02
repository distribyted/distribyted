Handlebars.registerHelper('ibytes', function (bytesSec, timePassed) {
    return Humanize.ibytes(bytesSec / timePassed, 1024);
});
Handlebars.registerHelper('bytes', function (bytes) {
    return Humanize.bytes(bytes, 1024);
});


var Distribyted = Distribyted || {};

Distribyted.message = {

    _toastr: function () {
        toastr.options = {
            closeButton: true,
            debug: false,
            newestOnTop: false,
            progressBar: true,
            positionClass: "toast-top-right",
            preventDuplicates: false,
            onclick: null,
            showDuration: "300",
            hideDuration: "1000",
            timeOut: "5000",
            extendedTimeOut: "1000",
            showEasing: "swing",
            hideEasing: "linear",
            showMethod: "fadeIn",
            hideMethod: "fadeOut"
        };

        return toastr;
    },


    error: function (message) {
        this._toastr().error(message);
    },

    info: function (message) {
        this._toastr().info(message);
    }
}

$(document).ready(function () {
    "use strict";

    /*======== 1. SCROLLBAR SIDEBAR ========*/
    var sidebarScrollbar = $(".sidebar-scrollbar");
    if (sidebarScrollbar.length != 0) {
        sidebarScrollbar.slimScroll({
            opacity: 0,
            height: "100%",
            color: "#808080",
            size: "5px",
            touchScrollStep: 50
        })
            .mouseover(function () {
                $(this)
                    .next(".slimScrollBar")
                    .css("opacity", 0.5);
            });
    }

    /*======== 2. MOBILE OVERLAY ========*/
    if ($(window).width() < 768) {
        $(".sidebar-toggle").on("click", function () {
            $("body").css("overflow", "hidden");
            $('body').prepend('<div class="mobile-sticky-body-overlay"></div>')
        });

        $(document).on("click", '.mobile-sticky-body-overlay', function (e) {
            $(this).remove();
            $("#body").removeClass("sidebar-mobile-in").addClass("sidebar-mobile-out");
            $("body").css("overflow", "auto");
        });
    }

    /*======== 3. SIDEBAR MENU ========*/
    var sidebar = $(".sidebar")
    if (sidebar.length != 0) {
        $(".sidebar .nav > .has-sub > a").click(function () {
            $(this).parent().siblings().removeClass('expand')
            $(this).parent().toggleClass('expand')
        })

        $(".sidebar .nav > .has-sub .has-sub > a").click(function () {
            $(this).parent().toggleClass('expand')
        })
    }


    /*======== 4. SIDEBAR TOGGLE FOR MOBILE ========*/
    if ($(window).width() < 768) {
        $(document).on("click", ".sidebar-toggle", function (e) {
            e.preventDefault();
            var min = "sidebar-mobile-in",
                min_out = "sidebar-mobile-out",
                body = "#body";
            $(body).hasClass(min)
                ? $(body)
                    .removeClass(min)
                    .addClass(min_out)
                : $(body)
                    .addClass(min)
                    .removeClass(min_out)
        });
    }

    /*======== 5. SIDEBAR TOGGLE FOR VARIOUS SIDEBAR LAYOUT ========*/
    var body = $("#body");
    if ($(window).width() >= 768) {

        if (typeof window.isMinified === "undefined") {
            window.isMinified = false;
        }
        if (typeof window.isCollapsed === "undefined") {
            window.isCollapsed = false;
        }

        $("#sidebar-toggler").on("click", function () {
            if (
                body.hasClass("sidebar-fixed-offcanvas") ||
                body.hasClass("sidebar-static-offcanvas")
            ) {
                $(this)
                    .addClass("sidebar-offcanvas-toggle")
                    .removeClass("sidebar-toggle");
                if (window.isCollapsed === false) {
                    body.addClass("sidebar-collapse");
                    window.isCollapsed = true;
                    window.isMinified = false;
                } else {
                    body.removeClass("sidebar-collapse");
                    body.addClass("sidebar-collapse-out");
                    setTimeout(function () {
                        body.removeClass("sidebar-collapse-out");
                    }, 300);
                    window.isCollapsed = false;
                }
            }

            if (
                body.hasClass("sidebar-fixed") ||
                body.hasClass("sidebar-static")
            ) {
                $(this)
                    .addClass("sidebar-toggle")
                    .removeClass("sidebar-offcanvas-toggle");
                if (window.isMinified === false) {
                    body
                        .removeClass("sidebar-collapse sidebar-minified-out")
                        .addClass("sidebar-minified");
                    window.isMinified = true;
                    window.isCollapsed = false;
                } else {
                    body.removeClass("sidebar-minified");
                    body.addClass("sidebar-minified-out");
                    window.isMinified = false;
                }
            }
        });
    }

    if ($(window).width() >= 768 && $(window).width() < 992) {
        if (
            body.hasClass("sidebar-fixed") ||
            body.hasClass("sidebar-static")
        ) {
            body
                .removeClass("sidebar-collapse sidebar-minified-out")
                .addClass("sidebar-minified");
            window.isMinified = true;
        }
    }
});