htmx.on("htmx:beforeSwap", function (evt) {

  var incomingDOM = new DOMParser().parseFromString(evt.detail.xhr.response, "text/html");

  var path = "head *[data-page-specific]";
  document.querySelectorAll(path).forEach(function (e) {
    e.parentNode.removeChild(e);
  });
  incomingDOM.querySelectorAll(path).forEach(function (e) {
    document.head.appendChild(e);
  })

});
