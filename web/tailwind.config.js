/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,html,svelte}",
  ],
  theme: {
    extend: {},
  },
  safelist: [
    "text-red-600",
    "text-lime-400",
  ],
  plugins: [],
}

