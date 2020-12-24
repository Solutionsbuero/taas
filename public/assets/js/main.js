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


// HANDLERS

function on_turnout_click(id) {
	console.log("hoi" + id);
}

function on_turnout_over(id) {
	set_status("Change the position of turnout " + id + ".");
}

function on_forward_click() {
	console.log("forward");
	document.getElementById("train-speed").textContent = "-4";
}

function on_forward_over() {
	set_status("Increase speed forward by one");
}

function on_backward_click() {
	console.log("forward");
}

function on_backward_over() {
	set_status("Increase speed backward by one");
}

function on_mouse_out() {
	set_status("...");
}

function on_ready() {
	document.getElementById("train-speed").textContent = "4";
}

ready(on_ready());
