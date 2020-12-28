let speed = -3;

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
	// document.getElementById("train-1-speed").textContent = "-4";
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
	// document.getElementById("train-0-speed").textContent = "4";
}

ready(on_ready());
