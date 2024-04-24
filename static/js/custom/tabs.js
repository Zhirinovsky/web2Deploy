$( document ).ready(function () {
    var hash = document.location.hash;
    var prefix = "tab_";
    if (hash) {
        $('.nav-tabs a[href="nav-product"]').tab('hide');
        $('.nav-tabs a[href="'+hash.replace(prefix,"")+'"]').tab('show');
        $('.dropdown-menu').removeClass('show');
    }

    // Change hash for page-reload
    $('.nav-tabs a').on('shown', function (e) {
        window.location.hash = e.target.hash.replace("#", "#" + prefix);
    });
});