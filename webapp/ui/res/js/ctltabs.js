var ctltabs = {

LoadContent: function(td, url) {
  if(!stdlib.class.Contains(td, 'active')) {
    var tds = stdlib.dom.GetParentTag(td, 'tr').children;
    var i;
    for (i = 0; i < tds.length; i++) {
      var ctd = tds[i];
      if (ctd.classList.contains('active')) {
        ctd.classList.remove('active');
      }
    }
    td.classList.toggle("active");
    var container = stdlib.dom.GetParentTag(td, 'div');
    stdlib.dom.LoadInnerHTML(url + '?id=' + container.id + '&tab=' + td.id, container.id + '_content');
  }
}

}
