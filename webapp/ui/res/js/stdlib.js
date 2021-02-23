/*------------------------------------------------------------------------

  Developer:  Iris Fagerland
  Created:    02-OCT-2017

  Javascript library of usefull utility functions:

  DOM:
    stdlib.dom.Ajax(method, url)
    stdlib.dom.GetChildNumber(pElement) 
    stdlib.dom.GetParentTag(pElement, pTag)
    stdlib.dom.LoadInnerHTML(url, id)

  Error Handling:
    stdlib.error.Log(src, msg)

  Utilities:
    stdlib.IsObject(candidate)
    stdlib.Pad(str, len)
    stdlib.Trim(str)

  Class Manipulation:
    stdlib.class.Contains(pElement, pClassName)
    stdlib.class.Add(pElement, pClassName)
    stdlib.class.Remove(pElement, pClassName)
    stdlib.class.Swap(pElement, pClassNameOld, pClassNameNew)
    stdlib.class.GetArray(pElement)
    stdlib.class.GetList(pElement)

------------------------------------------------------------------------*/

var stdlib = {

IsObject: function(candidate) {
  var type = typeof candidate;
  if(type == 'object') { 
    return true;
  }
  return false;
},

Pad: function(str, len) {
  if(str.length < len) {
    return ('0000000000' + str).slice(-1 * len);
  }
  return str;
},

Trim: function(str) {
  str.replace(/^\s*/, '').replace(/\s*$/, '');
  return str;
},

class: {

  src: 'class.',

  Contains: function(pElement, pClassName) {
    try {
      if(!pElement || !pClassName) { 
        return false; 
      }
      var l_class = stdlib.Trim(pClassName);
      var l_elem = stdlib.IsObject(pElement) ? pElement : document.getElementById(pElement);
      if(!l_elem) { 
        return false; 
      }
      var l_classList = l_elem.className;
      if(!l_classList) { 
        return false; 
      } else {
        var arr = l_classList.split(' ');
        for(var i = 0; i < arr.length; i++) {
          if(l_class == stdlib.Trim(arr[i])) { 
            return true; 
          }
        }
        return false;
      }
      return false;
    } catch(err) { 
      stdlib.error.Log(stdlib.class.src + 'Contains: ' + err); 
    }
  },

  Add: function(pElement, pClassName) {
    try {
      if(!pElement || !pClassName || stdlib.class.Contains(pElement, pClassName)) { 
        return; 
      }
      var l_elem = stdlib.IsObject(pElement) ? pElement : document.getElementById(pElement);
      if(!l_elem) {
        return false;
      }
      var l_classList = stdlib.Trim(l_elem.className);
      var arr = l_classList.split(' ');
      arr.push(stdlib.Trim(pClassName));
      l_elem.className = arr.join(' ');
    } catch(err) { 
      stdlib.error.Log(stdlib.src + 'Add: ' + err); 
    }
  },

  Remove: function(pElement, pClassName) {
    try {
      if(!pElement || !pClassName || !stdlib.class.Contains(pElement, pClassName)) {
        return; 
      }
      var l_elem = stdlib.IsObject(pElement) ? pElement : document.getElementById(pElement);
      if(!l_elem) {
        return false;
      }
      var l_classList = stdlib.Trim(l_elem.className);
      var arr = l_classList.split(' ');
      var arrResults = new Array();
      var l_class = stdlib.Trim(pClassName);
      for(var i = 0; i < arr.length; i++) {
        if(l_class != stdlib.Trim(arr[i])) { 
          arrResults.push(stdlib.Trim(arr[i])); 
        }
      }
      l_elem.className = arrResults.join(' '); 
    } catch(err) { 
      stdlib.error.Log(stdlib.src + 'Remove: ' + err); 
    }
  },

  Swap: function(pElement, pClassNameOld, pClassNameNew) {
    try {
      if(!pElement || !document.getElementById(pElement)) { 
        return; 
      }
      if(pClassNameOld) { 
        Remove(pElement, pClassNameOld); 
      }
      if(pClassNameNew) { 
        Add(pElement, pClassNameNew); 
      }
    } catch(err) { 
      stdlib.error.Log(stdlib.src + 'Swap: ' + err); 
    }
  },

  GetArray: function(pElement) {
    try {
      if(!pElement || !document.getElementById(pElement)) { 
        return; 
      }
      var l_elem = document.getElementById(pElement);
      var l_classList = stdlib.Trim(l_elem.className);
      var arr = l_classList.split(' ');
      var arrResults = new Array();
      for(var i = 0; i < arr.length; i++) {
        if(stdlib.Trim(arr[i]) != '') { 
          arrResults.push(stdlib.Trim(arr[i])); 
        }
      }
      return arrResults;
    } catch(err) { 
      stdlib.error.Log(stdlib.src + 'GetArray: ' + err); 
    }
  },

  GetList: function(pElement) {
    try {
      var arr = getClassArray(pElement);
      return arr.join(' ');
    } catch(err) { 
      stdlib.error.Log(stdlib.src + 'GetList: ' + err); 
    }
  }

}, // class

dom: {

  Ajax: function(method, url) {
    try {
      var http = new XMLHttpRequest();
      http.open(method, url, true);
      http.send(null);
    } catch(err) {
      console.log('stdlib.dom.LoadInnerHTML: ' + err);
    }
  },

  GetChildNumber: function(pElement) {
    var parent = pElement.parentNode;
    var children = parent.childNodes;
    for(var i = 0; i < children.length; i++) {
      if(children[i] == pElement) {
        return i;
      }
    }
    return -1;
  },
 
  GetParentTag: function(pElement, pTag) {
    try {
      if(pElement.tagName.toUpperCase() == pTag.toUpperCase()) {
        return pElement;
      } else {
        return stdlib.dom.GetParentTag(pElement.parentNode, pTag);
      }
    } catch(err) {
      stdlib.error.Log(stdlib.dom.src + 'GetParentTag: ' + err); 
    }
  },

  LoadInnerHTML: function(url, id) {
    try {
      var http = new XMLHttpRequest();
      http.open('GET', url, true);
      http.onreadystatechange = function() {
        if((http.readyState == 4) && (http.responseText != '')) {
          document.getElementById(id).innerHTML = http.responseText;
        }
      }
      http.send(null);
    } catch(err) {
      console.log('stdlib.dom.LoadInnerHTML: ' + err);
    }
  }

} // dom

}; // stdlib
