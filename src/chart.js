(() => {
  function renderChart(chart, target) {
    const closestDivEl = target ?? document.currentScript.previousSibling;

    const chartData = closestDivEl.getAttribute("data-data");
    const chartLabels = closestDivEl.getAttribute("data-labels");
    const chartId = closestDivEl.getAttribute("data-id");

    const data = {
      labels: JSON.parse(chartLabels),
      datasets: [
        {
          data: JSON.parse(chartData),
          backgroundColor: ["#1D4ED8"],
          borderRadius: 8,
        },
      ],
    };

    if (chart) {
      chart.data.labels = JSON.parse(chartLabels);
      chart.data.datasets = [
        {
          ...chart.data.datasets[0],
          data: JSON.parse(chartData),
        },
      ];
      chart.update();
      return;
    }

    chart = new Chart(document.getElementById(`chart-${chartId}`), {
      type: "bar",
      data: data,
      options: {
        plugins: {
          legend: {
            display: false,
          },
          tooltip: {
            enabled: true,
            backgroundColor: "#0a0a0a",
            titleColor: "white",
            bodyColor: "white",
            borderColor: "#262626",
            borderWidth: 1,
            padding: 10,
            cornerRadius: 8,
            displayColors: false,
            callbacks: {
              label: function (context) {
                const value = context.raw;
                return `${value} new users`;
              },
            },
          },
        },
        scales: {
          y: {
            beginAtZero: true,
            max: 1000,
            ticks: {
              display: false,
            },
            border: {
              display: false,
            },
            grid: {
              color: "#262626",
            },
          },
          x: {
            ticks: {
              color: "#fff",
            },
            grid: {
              color: "transparent",
            },
          },
        },
      },
    });

    return {
      chartId,
      chart,
    };
  }

  const { chart, chartId } = renderChart();

  document.body.addEventListener("htmx:afterSettle", function (evt) {
    if (evt.detail.target.id === `chart-card-${chartId}`) {
      renderChart(chart, evt.detail.target);
    }
  });
})();
