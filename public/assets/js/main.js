let last_status = {};

// HELPERS

function ready(fn) {
  if (document.readyState != 'loading'){
    fn();
  } else {
    document.addEventListener('DOMContentLoaded', fn);
  }
}

function set_status(message) {
	document.getElementById("status").innerHTML = message;
}

function post_json(url, data) {
	let xhr = new XMLHttpRequest();
	xhr.open("POST", url, true);
	xhr.setRequestHeader("Content-Type", "application/json");
	let d = JSON.stringify(data);
	xhr.send(d);
}


// HANDLERS

function on_turnout_click(id) {
	console.log("turnout " + id + " change");
	post_json("/api/turnout/" + id + "/change", {})
}

function on_turnout_over(id) {
	set_status("Change the position of turnout " + id + ".");
}

function on_forward_click(id) {
	console.log("train " + id + " speed +1");
	post_json("/api/train/" + id + "/speed", {"speed_delta": 1})
}

function on_forward_over(id) {
	set_status("Increase speed forward by one");
}

function on_backward_click(id) {
	console.log("train " + id + " speed -1");
	post_json("/api/train/" + id + "/speed", {"speed_delta": -1})
}

function on_backward_over(id) {
	set_status("Increase speed backward by one");
}

function on_mouse_out() {
	set_status("...");
}

function on_ready() {
	let loc = window.location;
	let proto = "ws:";
	if (loc.protocol === "https:") {
		proto = "wss:";
	}
	let url = proto + "//" + loc.host + loc.pathname + "ws";
	ws = new WebSocket(url);
	ws.onopen = function() {
		console.log("connected to websocket at " + url);
	}
	ws.onmessage = on_frontend_ws;
}

function on_frontend_ws(evt) {
	let data = JSON.parse(evt.data);
	last_status = data;

	update_train(1, data.train_1_speed);
	update_train(2, data.train_2_speed);
	update_turnout(0, data.turnout_0_position);
	update_turnout(1, data.turnout_1_position);
	update_turnout(2, data.turnout_2_position);
	update_turnout(3, data.turnout_3_position);
	update_turnout(4, data.turnout_4_position);
}

// UPDATE UI

function update_train(id, speed) {
	let ele = document.getElementById("train-" + id + "-speed");
	if (!ele) {
		return
	}
	ele.textContent = speed;
}

function update_turnout(id, position) {
	let straight = document.getElementById("tu-" + id + "-s");
	let diverging = document.getElementById("tu-" + id + "-d");
	if (!straight || !diverging) {
		return
	}

	if (position === 1) {
		straight.style.display = "block";
		diverging.style.display = "none";
	} else {
		straight.style.display = "none";
		diverging.style.display = "block";
	}
}

ready(on_ready());
