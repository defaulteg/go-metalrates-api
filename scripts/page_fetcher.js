var page = new WebPage()
var fs = require('fs');
var system = require('system');
var args = system.args;
var pathToScript, sourcePageAddress;

if (args.length < 2) {
    console.log("Error.\nNot enough params")    //Failed to execute phantomjs script: no parameters; only script file name was specified.
    phantom.exit();        //if something wrong return 1 ????
} else {
    pathToScript = args[0];
    sourcePageAddress = args[1];

    page.open(sourcePageAddress, function() {
            page.evaluate(function() {
        });
    });
}

page.onLoadFinished = function() {
    fs.write("./src/gitlab.com/defaulteg/api/scripts/temp/temp_data.html", page.content, 'w');
    phantom.exit();
};




