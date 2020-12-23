function ready(fn) {
  if (document.readyState != 'loading'){
    fn();
  } else {
    document.addEventListener('DOMContentLoaded', fn);
  }
}

function registerEventListeners() {
	document.addEventListener('click', function (event) {
		// if (!event.target.closest('.click-me')) return;
		console.log(event.target);
	}, false);
}

function on_turnout_click(id) {
	console.log("hoi" + id)
}

ready(registerEventListeners())
