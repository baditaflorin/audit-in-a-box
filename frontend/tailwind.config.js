/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        ink: "#17211b",
        paper: "#f7f5ee",
        moss: "#0f675f",
        coral: "#c44939",
        gold: "#c18a22",
        sky: "#256d9c",
      },
      boxShadow: {
        soft: "0 18px 45px rgba(23, 33, 27, 0.12)",
      },
    },
  },
  plugins: [],
};
