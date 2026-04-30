const stylelint = require('stylelint');

const colorHexUppercaseRuleName = 'icw/color-hex-uppercase';
const hexColorPattern = /#(?:[0-9a-fA-F]{3}|[0-9a-fA-F]{4}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})(?![0-9a-fA-F])/g;

const colorHexUppercaseMessages = stylelint.utils.ruleMessages(colorHexUppercaseRuleName, {
  rejected: (color) => `Expected "${color}" to use uppercase hex color notation.`,
});

const colorHexUppercaseRule = stylelint.createPlugin(colorHexUppercaseRuleName, (enabled, _options, context) => {
  return (root, result) => {
    const validOptions = stylelint.utils.validateOptions(result, colorHexUppercaseRuleName, {
      actual: enabled,
      possible: [true],
    });

    if (!validOptions || !enabled) {
      return;
    }

    root.walkDecls((decl) => {
      if (context.fix) {
        decl.value = decl.value.replace(hexColorPattern, (color) => color.toUpperCase());
        return;
      }

      const colors = decl.value.matchAll(hexColorPattern);
      for (const color of colors) {
        const rawColor = color[0];
        if (rawColor !== rawColor.toUpperCase()) {
          stylelint.utils.report({
            message: colorHexUppercaseMessages.rejected(rawColor),
            node: decl,
            result,
            ruleName: colorHexUppercaseRuleName,
          });
        }
      }
    });
  };
});

module.exports = {
  plugins: [colorHexUppercaseRule],
  rules: {
    'color-hex-length': 'long',
    'color-named': 'never',
    'function-disallowed-list': ['rgb', 'rgba', 'hsl', 'hsla'],
    [colorHexUppercaseRuleName]: true,
  },
};
