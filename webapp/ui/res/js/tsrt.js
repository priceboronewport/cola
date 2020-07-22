/*------------------------------------------------------------------------

  Developer:  Brian Fagerland
  Created:    02-OCT-2017

  Javascript to enable column sorting on tables.

  Requires:
    stdlib.js

  CSS:
    table th.sa {} // sorted ascending style
    table th.sd {} // sorted descending style

  Methods:
    tsrt.Sort(th_element)

  Usage:
    <th onClick='tsrt.Sort(this)'>Column Name</th>

------------------------------------------------------------------------*/

var tsrt = {

  AsDate: function(value) {
    var d = new Date(value);
    if(d != 'Invalid Date') {
      return stdlib.Pad(d.getFullYear()+'', 4) + stdlib.Pad(d.getMonth()+'', 2)
        + stdlib.Pad(d.getDate()+'', 2) + stdlib.Pad(d.getHours()+'', 2)
        + stdlib.Pad(d.getMinutes()+'', 2) + stdlib.Pad(d.getSeconds()+'', 2);
    } else {
      return '99999999999999';
    }
  },

  AsNumber: function(value) {
    var tmp = '';
    for(var i = 0; i < value.length; i++) {
      if(((value.substring(i, i + 1) >= '0') && (value.substring(i, i + 1) <= '9')) 
          || (value.substring(i, i + 1) == '-') || (value.substring(i, i + 1) == '.')) {
        tmp = tmp + value.substring(i, i + 1);
      }
    }
    if(tmp == '') {
      tmp = '0';
    }
    return parseFloat(tmp);
  },

  ClearSorts: function(thead) {
    var children = thead.childNodes;
    for(i in children) {
      stdlib.class.Remove(children[i], 'sa'); 
      stdlib.class.Remove(children[i], 'sd'); 
    }
  },

  Reverse: function(table) {
    var trs = table.tBodies[0].rows;
    for(var i = 0; i < trs.length; i++) {
      table.tBodies[0].insertBefore(trs[i], trs[0]);
    }
  },

  Sort: function(th) {
    try {
      var sort_type = 'text';
      if(stdlib.class.Contains(th, 'number')) {
        sort_type = 'number';
      } else if(stdlib.class.Contains(th, 'date')) {
        sort_type = 'date';
      }
      var table = stdlib.dom.GetParentTag(th, 'table');
      if(stdlib.class.Contains(th, 'sa')) {
        tsrt.Reverse(table);
        stdlib.class.Remove(th, 'sa');
        stdlib.class.Add(th, 'sd');
      } else if(stdlib.class.Contains(th, 'sd')) {
        tsrt.Reverse(table);
        stdlib.class.Remove(th, 'sd');
        stdlib.class.Add(th, 'sa');
      } else {
        var values = new Array();
        var col = stdlib.dom.GetChildNumber(th);
        var tbody = table.tBodies[0]
        for(var i = 0; i < tbody.rows.length; i++) {
          values.push(tbody.rows[i].cells[col].innerHTML);
        }
        var sorted = values.slice(0);
        if(sort_type == 'number') {
          sorted.sort(function(a,b) {
            var na = tsrt.AsNumber(a);
            var nb = tsrt.AsNumber(b);
            if(na > nb) {
              return 1;
            } else if(na < nb) {
              return -1;
            }
            return 0;
          });
        } else if(sort_type == 'date') {
          sorted.sort(function(a,b) {
            var da = tsrt.AsDate(a);
            var db = tsrt.AsDate(b);
            if(da > db) {
              return 1;
            } else if(da < db) {
              return -1;
            } 
            return 0;
          });
        } else {
          sorted.sort();
        }
        for(var i = 0; i < sorted.length; i++) {
          for(var j = 0; j < values.length; j++)
            if(sorted[i] == values[j]) {
              values[j] = values[i];
              values[i] = '~';
              if(i != j) {
                tsrt.SwapRows(table, i, j);
              }
              break;
            }
        }
        tsrt.ClearSorts(th.parentNode);
        stdlib.class.Add(th, 'sa');
      }
    } catch(err) {
      console.log('tsrt.Sort: ' + err);
    }
  },

  SwapRows: function(table, i, j) {
    try {
      var trs = table.tBodies[0].getElementsByTagName('TR');
      if(i == j+1) {
        table.tBodies[0].insertBefore(trs[i], trs[j]);
      } else if(j == i+1) {
        table.tBodies[0].insertBefore(trs[j], trs[i]);
      } else {
        var tmpNode = table.tBodies[0].replaceChild(trs[i], trs[j]);
        if(typeof(trs[i]) != 'undefined') {
          table.tBodies[0].insertBefore(tmpNode, trs[i]);
        } else {
          table.appendChild(tmpNode);
        }
      }  
    } catch(err) {
      console.log('tsrt.SwapRows: ' + err);
    }
  }

}; // tsrt
