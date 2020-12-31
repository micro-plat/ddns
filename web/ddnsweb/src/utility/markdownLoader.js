const markdown = require('markdown-it');

function markLoader(src) {
    const md = markdown({
        html: true,
        typographer: true,
    });
    const html = md.render(src);

    return (
        '<template>\n'
        + `<div class="markdown">${html}</div>\n`
        + '</template>\n'
    );
}
module.exports = markLoader;