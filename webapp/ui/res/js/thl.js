/*------------------------------------------------------------------------

  Developer:  Brian Fagerland
  Created:    02-OCT-2017

  Javascript to enable row highlighting within tables.

  Requires:
    stdlib.js

  CSS:
    table.multi_hl // allows multiple line highlighting
    table tbody tr.highlight {} // highlight style

  Methods:
    thl.Toggle(tr_element)

------------------------------------------------------------------------*/

var thl = {
  Toggle: function(tr) {
    try {
      var tbody = stdlib.dom.GetParentTag(tr, 'tbody');
      var table = stdlib.dom.GetParentTag(tr, 'table');
      if(!stdlib.class.Contains(table, 'multi_hl')) {
        if(stdlib.class.Contains(tr, 'highlight')) {
          stdlib.class.Remove(tr, 'highlight');
        } else {
          var children = tbody.childNodes;
          for(child in children) {
            stdlib.class.Remove(children[child], 'highlight'); 
          }
          stdlib.class.Add(tr, 'highlight');
        }
      } else {
        if(stdlib.class.Contains(tr, 'highlight')) {
          stdlib.class.Remove(tr, 'highlight');
        } else {
          stdlib.class.Add(tr, 'highlight');
        }
      }
    } catch(err) {
      console.log('thl.Toggle: ' + err);
    }
  }

}; // thl
