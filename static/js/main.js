Handlebars.registerHelper('prettyDate', function(unixDate) {
    if (unixDate == 0 || isNaN(unixDate) || typeof unixDate === undefined){
        return "";
    }
    var months = "Jan,Feb,Mar,Apr,May,Jun,Jul,Aug,Sep,Oct,Nov,Dec";
    function nth(d) {
      if(d>3 && d<21) return 'th'; // thanks kennebec
      switch (d % 10) {
            case 1:  return "st";
            case 2:  return "nd";
            case 3:  return "rd";
            default: return "th";
        }
    }
    d = new Date(1000 * unixDate);
    return d.getDate() + nth(d.getDate()) + " " + months.split(",")[d.getMonth()] + ", " + d.getFullYear() + " " + d.getHours() + ":" + d.getMinutes() + ":" + d.getSeconds();
});
// for stupid IE
$.ajaxSetup({ cache: false });
var app = {
    init : function(){
        $('select').material_select();
    }
}