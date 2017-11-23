var page = new WebPage()

var system = require('system');
var args = system.args;

if (args.length < 3) {
    console.log("Error.\nNot enough params");
    phantom.exit();
}

var targetSite = args[1]; //where to search rates;  args[0] = path to script


page.open(targetSite, function(status) {
  if (status !== 'success') {
        console.log('Error.\nUnable to access network');
  } else {
        // lists all values of property by id with /n delimiter
        for(var index = 2; index < args.length; index++) {
            var elementValue = page.evaluate(function(arg) {
                var element = document.getElementById(arg);
                if (element === null) {
                    return -1;
                } else {
                    return element.textContent;
                }
             }, args[index]);
             console.log(elementValue);
        }
  }
  phantom.exit();
});







