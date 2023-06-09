/** @type {import('tailwindcss').Config} */
const colors = require('tailwindcss/colors');

module.exports = {
  darkMode: 'class',
  content: ["**/*.html"],
  theme: {
	  extend: {
		  colors: {
			  primary: colors.blue,
			  base: colors.slate
		  },
		  fontFamily: {
			  satoshi: ['Satoshi', 'sans-serif'],
			  cabinet: ['CabinetGrotesk', 'sans-serif']
		  }
	  }
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ],
}

