/*------------------------------------------------------------------------
  Developer:    Iris Fagerland
  Created:      06-OCT-2017
  Updated:      23-OCT-2017
  Purpose:      Column Filtering.
  Dependency:   stdlib.js
  Assumptions:  Assumes all columns being filtered are in the 1st tBody
  Rules:        ! is a reserved character and indicates "NOT LIKE"
                in a multiple value search (separated by commas) any ! will negate all options
  Methods:
    tfltr.Apply(button)
    tfltr.Reset(button)
    tfltr.Filter(input)
    tfltr.HideRows(pTable, pColIndx, pMatch)
    tfltr.Store(name, value)

  Rules:
    wild card = * (0, 1 or more char)
    not equal = !
    is null = * (with no other search string)
    is not null !* (with no other search string)

------------------------------------------------------------------------*/

var tfltr = {

  Apply: function(button) {
    try {
      if (!button) {
        return false;
      } else {
        var tr = stdlib.dom.GetParentTag(button, 'TR');
        var table = stdlib.dom.GetParentTag(tr, 'TABLE');
        var tds = tr.getElementsByTagName('TD');
        for (var i = 0; i < tds.length; i++) {
          if (stdlib.class.Contains(tds[i], 'filter')) {
            if(tds[i].childNodes.length > 0 && tds[i].childNodes[0].tagName == 'INPUT') {
              tfltr.Filter(tds[i].childNodes[0]);
              tfltr.Store('tfltr.' + table.id + '.' + i, tds[i].childNodes[0].value);
            }
          }
        }
      }
    } catch(e) { 
      console.log(e);
    }
    return false;
  }, // Apply

  Reset: function(button) {
    try {
      if (!button) {
        return false;
      } else {
        var tr = stdlib.dom.GetParentTag(button, 'TR'); 
        var table = stdlib.dom.GetParentTag(tr, 'TABLE');
        var tds = tr.getElementsByTagName('TD');
        for (var i = 0; i < tds.length; i++) {
          if (stdlib.class.Contains(tds[i], 'filter')) {
            if((tds[i].childNodes.length > 0) && (tds[i].childNodes[0].tagName == 'INPUT')) {
              tds[i].childNodes[0].value = '';
              tfltr.Store('tfltr.' + table.id + '.' + i, '');
            }
          }
        }
      }
      var trs = stdlib.dom.GetParentTag(button, 'TABLE').tBodies[0].rows;
      for (var i = 0; i < trs.length; i++) {
        trs[i].style.display = ''
      }
    } catch(e) {
      console.log(e);
    }
    return false;
  }, // Reset

  Filter: function(input) {
    var delimiter = ',';
    var search = input.value;
    var td = stdlib.dom.GetParentTag(input, 'TD'); 
    var table = stdlib.dom.GetParentTag(td, 'TABLE');
    if (search.match(delimiter)) { 
      var array = new Array();
      var parsed = search.split(delimiter);
      for (var i = 0; i < parsed.length; i++) {
        array.push(parsed[i].trim());
      }
      tfltr.HideRows(table, td.cellIndex, array);
    } else { 
      if(search) { 
        tfltr.HideRows(table, td.cellIndex, search.trim());
      }
    }
  }, // Filter

  HideRows: function(pTable,pColIndx,pMatch) {
    //Note: Array will not honor match_types: not-null, is-null.
    //      These are only valid when used by a single string
    //
    var l_isArray = stdlib.IsObject(pMatch) ? true : false; //check to see if array
    var l_rows = pTable.tBodies[0].rows;
    var l_value = '';
    var l_cells = [];
    var l_match_type = 'equal'; //valid options: equal|not-equal|not-null|is-null|wild-equal|wild-not-equal

    //test to see if search is for a:
    //  match
    //  negated match (does not match)
    //  not null
    //  is null
    //  wildcard
    if ( pMatch.toString().match(/\*/) ){
      if (pMatch.toString() == '*') {
        l_match_type = 'not-null';
      }else{ //more than 1 char long
        if (pMatch.toString().match(/\!/)){
          if(pMatch.toString().length <= 2){
            l_match_type = 'is-null';
          }else{
            l_match_type = 'wild-not-equal';
          }
        }else{
          l_match_type = 'wild-equal';
        }
      }
    }else{ //non-wildcard
      if (pMatch.toString().match(/\!/)){
        l_match_type = 'not-equal';
      }else{
        l_match_type = 'equal';
      }
    }

    if(l_isArray){//check for all values in the array
      var l_isMatch = 0;
      for (var i=0;i<l_rows.length;i++) {
        l_cells = l_rows[i].cells;
        l_value = l_cells[pColIndx].innerHTML;
        if (l_match_type == 'not-equal' || l_match_type == 'wild-not-equal'){
          l_isMatch = 1;
        }else{l_isMatch = 0;}

        for (var j=0;j<pMatch.length;j++) {
          if(l_match_type == 'wild-not-equal'){
            if(l_value.toUpperCase().match(pMatch[j].replace(/[\!]/g,'').replace(/[\*]/g,'.*').toUpperCase())){
              l_isMatch = 0;
            }
          }else if(l_match_type == 'wild-equal'){
            if(l_value.toUpperCase().match(pMatch[j].replace(/[\*]/g,'.*').toUpperCase())){
              l_isMatch = 1;
            }
          }else if(l_match_type == 'not-equal'){
            if(l_value.toUpperCase().match(pMatch[j].replace(/[\!]/g,'').toUpperCase())){
              l_isMatch = 0;
            }
          }else if(l_match_type == 'equal'){
            if(l_value.toUpperCase().match(pMatch[j].toUpperCase())){
              l_isMatch = 1;
            }
          }
        }//end j loop
        if (l_isMatch==0){
          l_rows[i].style.display = 'none';
        }
      }//end i loop
    }else{ //not an array. Single string match

      for (var i=0;i<l_rows.length;i++) {
        l_cells = l_rows[i].cells;
        l_value = l_cells[pColIndx].innerHTML;

        if (l_match_type == 'not-null'){
          if(l_value.length == 0){
            l_rows[i].style.display = 'none';
          }
        }else if(l_match_type == 'is-null'){
          if(!l_value.length == 0){
            l_rows[i].style.display = 'none';
          }
        }else if(l_match_type == 'wild-not-equal'){
          if(l_value.toUpperCase().match(pMatch.replace(/[\!]/g,'').replace(/[\*]/g,'.*').toUpperCase())){
            l_rows[i].style.display = 'none';
          }
        }else if(l_match_type == 'wild-equal'){
          if(!l_value.toUpperCase().match(pMatch.replace(/[\*]/g,'.*').toUpperCase())){
            l_rows[i].style.display = 'none';
          }
        }else if(l_match_type == 'not-equal'){
          if(l_value.toUpperCase().match(pMatch.replace(/[\!]/g,'').toUpperCase())){
            l_rows[i].style.display = 'none';
          }
        }else if(l_match_type == 'equal'){
          if(!l_value.toUpperCase().match(pMatch.toUpperCase())){
            l_rows[i].style.display = 'none';
          }
        }
      }
    } //end array/string check
  }, // HideRows

  Store: function(name, value) {
    console.log('name='+name+' value='+value);
  }, // Store

  KeyPress: function(input, event) {
    if(event.keyCode === 13) {
      tfltr.Apply(input);
      return false;
    }
    return true;
  }, // KeyPress

}; //tfltr
