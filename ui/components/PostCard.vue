<template>
    <div class="card mb-3">
        <div class="row g-0">
            <div class="col-12 px-5 pt-3">
                <h3 class="mb-0 fs-5">Blog TitleBlog TitleBlog TitleBlog Title</h3>
                <hr class="my-hr" style="width: 90%;">
            </div>

            <span type="button" class="head-adge">置顶</span>

            <div class="create-date">
                <p class="fs-4 mb-0 text-center">29</p>
                <div class="year-month"><span class="pe-2">08月</span><span>2023</span></div>
            </div>

            <div class="col-md-12">
                <div class="card-body">
                    <div class="clearfix">
                        <img src="~/assets/images/gopher_cart.png"
                            class="img-fluid rounded-start float-start img-thumbnail me-3"
                            style="max-width: 200px; max-height: 160px;">
                        <p class="card-text text-muted" style="letter-spacing: 1.2px;">
                            通过 Bootstrap 所提供的 .img-fluid 类让图片支持响应式布局。其原理是将 max-width: 100%; 和 height: auto;
                            赋予图片，以便随父元素的宽度变化一起缩放。 通过 Bootstrap 所提供的 .img-fluid 类让图片支持响应式布局。其原理是将 max-width: 100%; 和 height:
                            auto; 赋予图片，以便随父元素的宽度变化一起缩放。 通过 Bootstrap 所提供的 .img-fluid 类让图片支持响...
                        </p>
                    </div>

                    <div class="d-flex flex-md-wrap position-relative" style="height:40px; line-height: 40px;">
                        <hr class="flex-grow-1 my-hr" style="margin-top: 20px;">
                        <NuxtLink to="/posts/100">
                            <button type="button" class="btn btn-link" style="text-decoration: none;">继续阅读&gt;&gt;</button>
                        </NuxtLink>
                    </div>
                    <footer class="d-flex justify-content-between">
                        <div class="tags">
                            <Icon name="pajamas:label" color="orange" size="22" />
                            <button type="button" class="btn btn-primary btn-sm ms-2" style="--bs-btn-padding-y:0rem;">Golang</button>
                            <button type="button" class="btn btn-primary btn-sm ms-2" style="--bs-btn-padding-y:0rem;">Vue</button>
                        </div>
                        <div class="tool">
                            <button class="btn btn-light btn-sm ms-2 text-muted d-inline-flex align-items-center me-3" disabled style="--bs-btn-disabled-opacity: .9">
                                <Icon name="icon-park-outline:eyes" size="18" /> <span class="ms-1">200</span>
                            </button>
                            <button class="btn btn-light btn-sm ms-2 text-muted d-inline-flex align-items-center me-3">
                                <Icon name="icon-park-outline:thumbs-up" size="18" /> <span class="ms-1">200</span>
                            </button>
                            <button class="btn btn-light btn-sm ms-2 text-muted d-inline-flex align-items-center me-3">
                                <Icon name="icon-park-outline:copy" size="18" /> <span class="ms-1">分享</span>
                            </button>                            
                        </div>
                    </footer>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">

</script>

<style scoped lang="scss">
.card {
    position: relative;
    overflow: hidden;

    .head-adge {
        position: absolute;
        display: inline-block;
        top: 0;
        width: 120px;
        height: 26px;
        line-height: 26px;
        left: -34px;
        font-size: 12px;
        color: #ffffff;
        text-indent: 35px;
        background: #2470dc;
        transform: rotate(310deg);
    }

    .create-date {
        position: absolute;
        width: auto;
        top: 0;
        right: 0;
        padding: 5px 20px 0 12px;
    }

    .my-hr {
        border-top-style: dashed;
        border-width: 1.4px;
        border-color: var(--bs-gray-dark);
        margin-top: 10px;
        margin-bottom: 0;
    }
}
</style>