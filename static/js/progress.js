(function () {
  "use strict";

  var bar = document.getElementById("progress-bar");

  if (!bar) {
    // Progress bar element not found.
    return;
  }

  var completeTimer = null;

  // htmx:beforeRequest -> begin progress animation.
  document.addEventListener("htmx:beforeRequest", function () {
    if (completeTimer) {
      clearTimeout(completeTimer);
      completeTimer = null;
    }
    bar.classList.remove("complete");
    // Force reflow so remove/add class changes are seen as separate frames.
    void bar.offsetWidth;
    bar.classList.add("loading");
  });

  // htmx:afterRequest -> finish then fade out.
  document.addEventListener("htmx:afterRequest", function () {
    bar.classList.remove("loading");
    bar.classList.add("complete");
    completeTimer = setTimeout(function () {
      bar.classList.remove("complete");
      completeTimer = null;
    }, 400);
  });
})();
