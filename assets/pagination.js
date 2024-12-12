for (const el of document.querySelectorAll("#pagination button")) {
  el.addEventListener("click", function () {
    spinnerEl.classList.remove("hidden");
    spinnerEl.classList.add("block");

    tableEl.classList.add("hidden");
    tableEl.classList.remove("block");
  });
}
