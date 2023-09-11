<template>
    <div class="container-fulid p-3 position-relative">
        <button class="btn btn-primary position-fixed pulish" data-bs-toggle="offcanvas" data-bs-target="#offcanvasRight"
            aria-controls="offcanvasRight">发布</button>
        <div class="row">
            <div class="col-6">
                <textarea class="form-control post-content" placeholder="写文章..." v-model="content"
                    @input="watchInput()"></textarea>
            </div>
            <div class="col-6">
                <div class="content" v-html="markdown.postDom"></div>
            </div>
        </div>
    </div>

    <div class="offcanvas offcanvas-end" tabindex="-1" id="offcanvasRight" aria-labelledby="offcanvasRightLabel"
        style="--bs-offcanvas-width:500px">
        <div class="offcanvas-header">
            <h5 id="offcanvasRightLabel">发布文章</h5>
            <button type="button" class="btn-close text-reset" data-bs-dismiss="offcanvas" aria-label="Close"></button>
        </div>
        <div class="offcanvas-body">
            <form>
                <div class="mb-3">
                    <label for="title" class="form-label">标题</label>
                    <input type="type" class="form-control" id="title">
                </div>

                <div class="mb-3">
                    <label for="abstract" class="form-label">摘要</label>
                    <textarea type="text" class="form-control" id="abstract"></textarea>
                </div>

                <div class="input-group mb-3">
                    <label class="input-group-text" for="inputGroupFile01">Upload</label>
                    <input type="file" class="form-control" id="inputGroupFile01">
                </div>

                <div class="mb-3">
                    <label for="category" class="form-label">Category</label>
                    <VSelect id="category" />
                </div>

                <div class="mb-3">
                    <label for="tag" class="form-label">Tag</label>
                    <VSelect id="tag" />
                </div>

                <div class="mb-3">
                    <label for="badge" class="form-label">Badge</label>
                    <select class="form-select" aria-label="Default select example" id="badge">
                        <option selected>Open this select menu</option>
                        <option value="1">One</option>
                        <option value="2">Two</option>
                        <option value="3">Three</option>
                    </select>
                </div>


                <div class="mb-3">
                    <button type="submit" class="btn btn-primary">Submit</button>
                </div>
            </form>
        </div>
    </div>
</template>

<script lang="ts" setup>

definePageMeta({
    layout: "blank"
})


const content = ref<string>("")
const value = ref<string>("")
const options = ref([
    'Batman',
    'Robin',
    'Joker',
])

const { markdown, parseMD } = useMarkdownParse()

let timer: any;
const watchInput = () => {
    if (timer) {
        clearTimeout(timer)
    }
    timer = setTimeout(() => {
        parseMD(content.value)
    }, 200)
}

</script>


<style scoped>
@import "vue-select/dist/vue-select.css";
</style>

<style lang="scss" scoped>
@import url("highlight.js/scss/monokai-sublime.scss");

.pulish {
    right: 30px;
}

.post-content {
    width: 100%;
    height: 100%;
    min-height: 95vh;
}
</style>