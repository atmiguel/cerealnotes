$(function() {
  // http://stepansuvorov.com/blog/2014/04/jquery-put-and-delete/
  jQuery.each( [ "put", "delete" ], function( i, method ) {
    jQuery[ method ] = function( url, data, callback, type ) {
      if ( jQuery.isFunction( data ) ) {
        type = type || callback;
        callback = data;
        data = undefined;
      }
   
      return jQuery.ajax({
        url: url,
        type: method,
        dataType: type,
        data: data,
        success: callback
      });
    };
  });

  jQuery.prototype.getDOM = function() {
    if (this.length === 1) {
        return this[0];
    }

    if (this.length === 0) {
      throw "jQuery object is empty"
    }
    throw "jQuery Object contains more than 1 object";
  };

});