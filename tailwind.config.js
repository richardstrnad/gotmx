/** @type {import('tailwindcss').Config} */
const colors = require('tailwindcss/colors');
const plugin = require('tailwindcss/plugin');

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
    plugin(function({ addVariant }) {
      addVariant('htmx-settling', ['&.htmx-settling', '.htmx-settling &'])
      addVariant('htmx-request',  ['&.htmx-request',  '.htmx-request &'])
      addVariant('htmx-swapping', ['&.htmx-swapping', '.htmx-swapping &'])
      addVariant('htmx-added',    ['&.htmx-added',    '.htmx-added &'])
    }),
  ],
}

