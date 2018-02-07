function httpGet(theUrl) {
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.open("GET", theUrl, false); // false for synchronous request
    xmlHttp.send(null);
    return xmlHttp.responseText;
}
var js = httpGet("/list")
var files = JSON.parse(js)

for (var i = 0; i < files.length; i++) {
    var file = files[i]
    var operations = {
        "Resize": "width=300&height=300",
        "Crop": "crop=150,150|250,250",
        "Rotate": "rotation=90&background=00FF00"
    };
    var container = document.getElementById('results');
    var title = document.createElement('div');

    title.innerHTML = "File name: <b>" + file + "</b>";

    container.appendChild(title);
    var table = document.createElement('table');
    table.className = 'table table-bordered';
    var thead = document.createElement('thead');
    var th1 = document.createElement('th');
    th1.innerText = 'Operations';
    thead.appendChild(th1);
    for (var key in operations) {
        var th = document.createElement('th');
        th.innerText = key;
        thead.appendChild(th);
    }
    table.appendChild(thead);
    var tbody = document.createElement('tbody');

    for (var key in operations) {
        var tr = document.createElement('tr');
        var th = document.createElement('th');

        th.innerText = key;
        tr.appendChild(th);

        for (var key2 in operations) {
            var td = document.createElement('td');
            var op2 = key != key2 ? '&' + operations[key2] : '';
            var title2 = key != key2 ? ' and ' + key2 : '';
            var fileName = file.split('.')[0];
            var link = '/src?source=' + fileName + '&' + operations[key] + op2;
            td.innerHTML = '<a href="' + link + '">' + key + ' ' + ' ' + title2 + '</a>' +
                ' (<a href="' + link + '&format=jpg">JPG</a>,' +
                ' <a href="' + link + '&format=png">PNG</a>,' +
                ' <a href="' + link + '&format=gif">GIF</a>)';
            tr.appendChild(td);
        }
        tbody.appendChild(tr);
    }
    table.appendChild(tbody);
    container.appendChild(table);
}