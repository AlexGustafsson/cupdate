import autoprefixer from 'autoprefixer'
import tailwindcss from 'tailwindcss'

// Markdown can have imgs using the "height" attribute. This is overridden by
// tailwind's defaults as it has a higher specificity. The rule is unused by us
// and can be introduced manually when needed instead.
// SEE: https://github.com/tailwindlabs/tailwindcss/pull/7742#issuecomment-1061332148
const removeTailwindImgAutoHeight = (css) => {
  css.walkRules(function (rule) {
    if (rule.selector === 'img,\nvideo') {
      rule.walkDecls(function (decl) {
        if (['height'].indexOf(decl.prop) !== -1) {
          decl.remove()
        }
      })
    }
  })
}

export default {
  plugins: [tailwindcss, autoprefixer, removeTailwindImgAutoHeight],
}
