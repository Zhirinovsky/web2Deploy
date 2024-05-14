new purePajinate({
    itemsPerPage: $('#pageSelect').val(),
    wrapAround: true,
    navLabelFirst: '<i class="bi bi-skip-start-fill"></i>',
    navLabelPrev: '<i class="bi-caret-left-fill"></i>',
    navLabelNext: '<i class="bi-caret-right-fill"></i>',
    navLabelLast: '<i class="bi bi-skip-end-fill"></i>',
    navOrder: ["first", "prev", "num", "next", "last"],
    showFirstLast: true,
    showPrevNext: true
});