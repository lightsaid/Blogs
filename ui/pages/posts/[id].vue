<template>
    <div class="container-fulid posts">
        <div class="wrapper">
            <header>
                <h1 class="text-lg font-medium text-slate-900">Go 语言 15 个内置函数详解</h1>
                <footer class=" mt-3">
                    <p class="text-sm text-gray-500">
                        <span>2023-05-15</span>
                        <span> · 3 min · LightSaid</span>
                    </p>
                </footer>
            </header>

            <div class="fixed right-40 top-32">
                <div v-html="markdown.menuDom"></div>
            </div>
            <div v-html="markdown.postDom"></div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import "highlight.js/styles/monokai-sublime.css"
import { markdownText } from "@/constants/markdown"
import { ref } from "vue"

definePageMeta({
    layout: "default"
})

const { markdown, parseMD, } = useMarkdownParse()

parseMD(markdownText)

</script>

<style scoped lang="scss">
:root{
    --gap: 24px;
    --navbar-width: 1024px;
    --main-width: 768px;
}

.wrapper {
    position: relative;
    max-width: calc(var(--main-width) + var(--gap) * 2);
    margin: auto;
    padding: var(--gap);
}

code {
    max-height: 650px;
    font-size: 0.875rem;
    border-radius: 5px;
    overflow: auto;
}

</style>



