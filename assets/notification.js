document.body.addEventListener("htmx:wsOpen", function () {
  const pulseEl = document.getElementById("pulse");
  pulseEl.classList.add("active");
});

document.body.addEventListener("htmx:wsClose", function () {
  const pulseEl = document.getElementById("pulse");
  pulseEl.classList.remove("active");
});
