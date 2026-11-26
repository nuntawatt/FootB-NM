/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        sport: {
          bg: '#0a0a0a',
          card: '#1a1a1a',
          hover: '#2a2a2a',
          border: '#333333',
          red: '#dc2626',
          'red-hover': '#ef4444',
          'red-dark': '#b91c1c',
        },
      },
      fontFamily: {
        sans: ['Noto Sans Thai', 'Inter', 'system-ui', 'sans-serif'],
      },
    },
  },
  plugins: [],
}
