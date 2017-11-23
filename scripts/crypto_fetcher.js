var page = new WebPage()

var system = require('system');
var args = system.args;

var targetSite = args[1]; //where to search rates;  args[0] = path to script

page.open(targetSite, function() {
    var x = page.evaluate(function() {
        var temp =  document.getElementsByClassName("table no-border table-condensed");
        var table = temp[0];
        var rateCount = 10;

        // If table html table doesn't contain enough elements just parse all elements there are
        if (table.rows.length < rateCount) {
            rateCount = table.rows.length - 1;
        }

        var data = [];
        data[0] = table.rows[1].cells[2].textContent.split('/')[0];    //currency name. e.g.: BTC; DOGE
         for (var i = 1, j = 1; i <= rateCount; i++) {
             var market = table.rows[i].cells[1].textContent;
             var pair = table.rows[i].cells[2].textContent;
             var value = table.rows[i].cells[4].textContent.replace('$',' ').trim();

             data[j++] = market;
             data[j++] = pair;
             data[j++] = value;
         }

         return data;
    });

    console.log(x);

    phantom.exit();
});


page.onLoadFinished = function() {


};




