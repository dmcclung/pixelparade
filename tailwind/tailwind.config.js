/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["/templates/**/*.{gohtml,html}"],
  theme: {
    extend: {},
  },
  plugins: [
    require("/tailwind/node_modules/@tailwindcss/forms"),
  ],
}

