$(function() {
  $('a[data-toggle="tab"]').on('shown.bs.tab', function (e) {
    var tab = $(e.target);
    var oldTab = $(e.relatedTarget);

    var id = tab.attr('href').replace("#", "");
    var oldId = oldTab.attr('href').replace("#", "");

    $('.tab-outer').removeClass(oldId);
    $('.tab-outer').addClass(id);
  })
});