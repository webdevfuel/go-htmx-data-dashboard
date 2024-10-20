/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["view/**/*.{html,js,templ}"],
  theme: {
    extend: {},
  },
  plugins: [require("@tailwindcss/forms")],
  darkMode: "class",
};
