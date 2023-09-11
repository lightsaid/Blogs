export const  markdownText = `# Hello, world!
This is a simple blog post written in **markdown** format.

## Section 1
This is the first section of the post.
Lorem ipsum dolor sit amet consectetur adipisicing elit. Placeat consequatur eius modi libero dignissimos dolor delectus, voluptates cupiditate natus a, itaque aperiam inventore sequi quo nulla harum perferendis accusantium quae.
Lorem ipsum dolor sit amet consectetur adipisicing elit. Placeat consequatur eius modi libero dignissimos dolor delectus, voluptates cupiditate natus a, itaque aperiam inventore sequi quo nulla harum perferendis accusantium quae.
Lorem ipsum dolor sit amet consectetur adipisicing elit. Placeat consequatur eius modi libero dignissimos dolor delectus, voluptates cupiditate natus a, itaque aperiam inventore sequi quo nulla harum perferendis accusantium quae.
Lorem ipsum dolor sit amet consectetur adipisicing elit. Placeat consequatur eius modi libero dignissimos dolor delectus, voluptates cupiditate natus a, itaque aperiam inventore sequi quo nulla harum perferendis accusantium quae.

\`\`\`go 
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         app.config.HTTPServerAddress,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		q := <-quit

		log.Println("recve signal", q.String())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		log.Println("do background tasks")

		// 等待在后台执行的任务完成
		app.wg.Wait()
		shutdownError <- nil
	}()

	log.Println("server running on ", app.config.HTTPServerAddress)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	log.Println("stopped server.")

	return nil
}

\`\`\`


### Subsection 1.1
This is a subsection of section 1.
Lorem ipsum dolor sit amet consectetur adipisicing elit. Placeat consequatur eius modi libero dignissimos dolor delectus, voluptates cupiditate natus a, itaque aperiam inventore sequi quo nulla harum perferendis accusantium quae.
Lorem ipsum dolor sit amet consectetur adipisicing elit. Placeat consequatur eius modi libero dignissimos dolor delectus, voluptates cupiditate natus a, itaque aperiam inventore sequi quo nulla harum perferendis accusantium quae.
Lorem ipsum dolor sit amet consectetur adipisicing elit. Placeat consequatur eius modi libero dignissimos dolor delectus, voluptates cupiditate natus a, itaque aperiam inventore sequi quo nulla harum perferendis accusantium quae.

### Subsection 1.2
This is another subsection of section 1.
Lorem ipsum dolor sit amet consectetur adipisicing elit. Placeat consequatur eius modi libero dignissimos dolor delectus, voluptates cupiditate natus a, itaque aperiam inventore sequi quo nulla harum perferendis accusantium quae.
Lorem ipsum dolor sit amet consectetur adipisicing elit. Placeat consequatur eius modi libero dignissimos dolor delectus, voluptates cupiditate natus a, itaque aperiam inventore sequi quo nulla harum perferendis accusantium quae.
Lorem ipsum dolor sit amet consectetur adipisicing elit. Placeat consequatur eius modi libero dignissimos dolor delectus, voluptates cupiditate natus a, itaque aperiam inventore sequi quo nulla harum perferendis accusantium quae.

\`\`\`html
<template>
    <div class="posts">
        <div class="wrapper">
            <header>
                <h1 class="text-lg font-medium text-slate-900">Go 个内置函数详解个内置函数详解个内置函数详解个内置函数详解 语言 15 个内置函数详解</h1>
                <footer class=" mt-3">
                    <p class="text-sm text-gray-500">
                        <span>2023-05-15</span>
                        <span> · 3 min · LightSaid</span>
                    </p>
                </footer>
            </header>
            <div v-html="htmlPost"></div>
        </div>
    </div>
</template>
\`\`\` 



## Section 2
This is the second section of the post.
Lorem ipsum dolor sit amet consectetur adipisicing elit. Placeat consequatur eius modi libero dignissimos dolor delectus, voluptates cupiditate natus a, itaque aperiam inventore sequi quo nulla harum perferendis accusantium quae.
Lorem ipsum dolor sit amet consectetur adipisicing elit. Placeat consequatur eius modi libero dignissimos dolor delectus, voluptates cupiditate natus a, itaque aperiam inventore sequi quo nulla harum perferendis accusantium quae.
Lorem ipsum dolor sit amet consectetur adipisicing elit. Placeat consequatur eius modi libero dignissimos dolor delectus, voluptates cupiditate natus a, itaque aperiam inventore sequi quo nulla harum perferendis accusantium quae.
Lorem ipsum dolor sit amet consectetur adipisicing elit. Placeat consequatur eius modi libero dignissimos dolor delectus, voluptates cupiditate natus a, itaque aperiam inventore sequi quo nulla harum perferendis accusantium quae.
Lorem ipsum dolor sit amet consectetur adipisicing elit. Placeat consequatur eius modi libero dignissimos dolor delectus, voluptates cupiditate natus a, itaque aperiam inventore sequi quo nulla harum perferendis accusantium quae.
Lorem ipsum dolor sit amet consectetur adipisicing elit. Placeat consequatur eius modi libero dignissimos dolor delectus, voluptates cupiditate natus a, itaque aperiam inventore sequi quo nulla harum perferendis accusantium quae.

### Subsection 2.1
This is a subsection of section 2.

\`\`\` js
var renderer = new marked.Renderer();
var headings = [];
renderer.heading = function (text, level) {
    var anchor = 'heading-' + headings.length;
    headings.push(anchor);
    return '<h' + level + ' id="' + anchor + '">' + text + '</h' + level + '>';
};

marked.setOptions({
    renderer: renderer,
    highlight: function (code, lang) {
        // console.log(code)
        const language = hljs.getLanguage(lang) ? lang : 'plaintext';
        let value = hljs.highlight(code, { language }).value
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

htmlPost.value = marked.parse(markdownText)

onMounted(()=>{
    // document.querySelectorAll(".copyBtn")
    new ClipboardJS('.copyBtn');
})

\`\`\`

### Subsection 2.2
This is another subsection of section 2.`;
