var prevData = "";
var prefix = "";

var levelMap = {
	"OK": 0,
	"INFO": 1,
	"WARNING": 2,
	"CRITICAL": 3,
}

moment.locale("DE");
Handlebars.registerHelper("formatDate", function(datetime) {
  if (moment) {
    return moment(datetime).format('LLL');
  }
  else {
    return datetime;
  }
});

function alertSorter(a, b) {
	if (levelMap[a["level"]] < levelMap[b["level"]]) return 1;
	if (levelMap[a["level"]] > levelMap[b["level"]]) return -1;
	if (a["time"] < b["time"]) return 1;
	if (a["time"] < b["time"]) return 1;
}

var loadAlerts = function() {
	var request = new XMLHttpRequest();
	request.open("GET", "/api/alerts" + prefix, true);
	request.onload = function() {
		if (request.responseText !== prevData) {
			var source = document.getElementById("alert-template").innerHTML;
			var template = Handlebars.compile(source);
			var elements = JSON.parse(request.responseText);
			elements.sort(alertSorter);

			var data = {
				alerts: elements
			};

			document.getElementById("alert-placeholder").innerHTML = template(data);

			// color the rows based on the error level
			var alert_tr = document.querySelectorAll("#alert-placeholder > table > tbody > tr");
			var v = ""
			for (var i = 0; i < alert_tr.length; i++) {
				var c = alert_tr[i].cells[0];
				switch (c.innerHTML) {
					case "OK":
						c.parentElement.style.background = "#2E7D32";
						break
					case "INFO":
						c.parentElement.style.background = "#1565C0";
						break
					case "WARNING":
						c.parentElement.style.background = "#EF6C00";
						break
					case "CRITICAL":
						c.parentElement.style.background = "#B71C1C";
						break
				}
			}
		};
		prevData = request.responseText;
	};
	request.send();
}

setInterval(loadAlerts, 1000);

// strange dom ready code - juhu no jquery
r(function() {
	loadAlerts();
});

function r(f) {
	/in/.test(document.readyState) ? setTimeout('r(' + f + ')', 9) : f()
}

function deleteAlert(element) {
	var id = element.parentElement.parentElement.cells[1].innerText;
	var request = new XMLHttpRequest();
	var encoded = btoa(id);
	// Replace characters according to base64url specifications
	encoded = encoded.replace(/\+/g, '-');
	encoded = encoded.replace(/\//g, '_');

	request.open("DELETE", "/api/alert/" + encoded, true);
	request.onload = function() {
		if (JSON.parse(request.responseText).id == id) {
			element.parentElement.parentElement.remove();
		}
	}
	request.send();
}
