function is_change(input) {
  is_search('is?id='+input.id+'_list&q=' + input.value, input.id);
  var list = document.getElementById(input.id + '_list');
  list.style.display = 'inline';
}

function is_search(url, id) {
  var http = new XMLHttpRequest();
  http.open('GET', url, true);
  var list = document.getElementById(id + '_list');
  list.innerHTML = '';
  http.onreadystatechange = function() {
    if((http.readyState == 4) && (http.responseText != '')) {
      is_load_list(http.responseText, id);
    }
  }
  http.send(null);
}

function is_select(button, input_id) {
  var list = document.getElementById(input_id + '_list');
  var input = document.getElementById(input_id);
  input.value = button.innerHTML;
  input.focus();
  list.style.display = 'none';
  return false;
}

function is_blur(input) {
  setTimeout(function() {
    if(document.activeElement.tagName != 'BUTTON') {
      var list = document.getElementById(input.id + '_list');
      list.style.display = 'none';
    }
  }, 400);
}

function is_blur_button(button, id) {
  setTimeout(function() {
    if(document.activeElement.tagName != 'BUTTON') {
      var list = document.getElementById(id + '_list');
      list.style.display = 'none';
    }
  }, 400);
}

function is_load_list(str, id) {
  var str_list = str.split(',');
  var content = '';
  for(var i = 0; i < str_list.length; i++) {
    content += "<button onBlur='is_blur_button(this, \"" + id + "\")' onClick='return is_select(this, \"" + id + "\")'/>" + str_list[i] + "</button><br/>";
  }
  var list = document.getElementById(id + '_list');
  list.innerHTML = content;
  return false;
}
