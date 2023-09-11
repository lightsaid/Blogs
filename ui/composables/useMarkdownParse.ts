import { reactive } from "vue"
import { marked } from "marked"
import hljs from "highlight.js"

export type Markdown = {
    mdDoc: string;
    menuDom: string;
    postDom: string;
}

export const useMarkdownParse = () => {
    const markdown = reactive<Markdown>({
        mdDoc: "",
        menuDom: "",
        postDom: ""
    })

    const parseMD = (mdDoc: string) => {
        markdown.mdDoc = mdDoc

        var renderer = new marked.Renderer();
        var headings: string[] = [];

        renderer.heading = function (text, level) {
            var anchor = 'heading-' + headings.length;
            headings.push(anchor);
            return '<h' + level + ' id="' + anchor + '">' + text + '</h' + level + '>';
        };

        renderer.table = function (header, body) {
            // 构建HTML表格
            return `<table border="1"><thead>${header}</thead><tbody>${body}</tbody></table>`;
          }

        marked.setOptions({
            renderer: renderer,
            highlight: function (code, lang) {
                console.log("code->",code)
                const language = hljs.getLanguage(lang) ? lang : 'plaintext';
                let value = hljs.highlight(code, { language }).value
                console.log("->>>",value)
                value += "<div>" + value + "</div>"

                return value
            },
            langPrefix: 'hljs language-', // highlight.js css expects a top-level 'hljs' class.
            pedantic: false,
            gfm: true,
            breaks: false,
            sanitize: false,
            smartypants: false,
            xhtml: false
        })

        markdown.postDom = marked.parse(mdDoc)

        markdown.menuDom = '<ul>';
        for (var i = 0; i < headings.length; i++) {
             markdown.menuDom  += '<li><a href="#' + headings[i] + '">' + headings[i] + '</a></li>';
        }
        markdown.menuDom  += '</ul>';
    }

    return {
        markdown,
        parseMD,
    }
}