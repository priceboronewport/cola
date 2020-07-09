function UIMenuShow(content) {
  window.scrollTo(0,0);
  document.getElementById(content).classList.toggle("ui_menu_show");
}

window.onclick = function(event) {
  if (!event.target.matches('.ui_menu_button')) {
    var dropdowns = document.getElementsByClassName("ui_menu_content");
    var i;
    for (i = 0; i < dropdowns.length; i++) {
      var openDropdown = dropdowns[i];
      if (openDropdown.classList.contains('ui_menu_show')) {
        openDropdown.classList.remove('ui_menu_show');
      }
    }
  }
}
