const usersEl = document.getElementById("users");

function dispatchRefetchEvent() {
  const refetchEvent = new CustomEvent("refetch");
  usersEl.dispatchEvent(refetchEvent);
}

const sortEl = document.getElementById("sort");

sortEl.addEventListener("change", function () {
  dispatchRefetchEvent();
});

const filterEl = document.getElementById("filter");

filterEl.addEventListener("change", function () {
  dispatchRefetchEvent();
});

const queryEl = document.getElementById("query");

let timeout;

queryEl.addEventListener("input", function () {
  clearTimeout(timeout);
  timeout = setTimeout(() => {
    dispatchRefetchEvent();
  }, 500);
});

const spinnerEl = document.getElementById("spinner");

const tableEl = document.getElementById("table");

document.body.addEventListener("htmx:afterRequest", function (evt) {
  if (evt.detail.target === tableEl) {
    spinnerEl.classList.add("hidden");
    spinnerEl.classList.remove("block");

    tableEl.classList.remove("hidden");
    tableEl.classList.add("block");
  }
});
