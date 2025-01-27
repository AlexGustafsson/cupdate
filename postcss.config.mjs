// Markdown can have imgs using the "height" attribute. This is overridden by
// tailwind's defaults as it has a higher specificity. The rule is unused by us
// and can be introduced manually when needed instead.
// SEE: https://github.com/tailwindlabs/tailwindcss/pull/7742#issuecomment-1061332148
const removeTailwindImgAutoHeight = (css) => {
  css.walkRules((rule) => {
    if (rule.selector === 'img, video') {
      rule.walkDecls((decl) => {
        if (['height'].indexOf(decl.prop) !== -1) {
          decl.remove()
        }
      })
    }
  })
}

/** @type {import('postcss-load-config').Config} */
const config = {
  plugins: [removeTailwindImgAutoHeight],
}

export default config
