<template>
    <aside class="card sticky-lg-top shadow" style="top: 83px; overflow: hidden; border: none; background-color: white;">
        <div class="position-sticky pt-4 pb-2">
            <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mb-1 text-muted">
                <div class="input-group mb-2 btn-search">
                    <input type="text" class="form-control rounded-pill my-input" placeholder="输入搜索内容">
                    <div class="position-absolute" style="right: 20px; top: 8px; z-index: 99;">
                        <Icon name="simple-line-icons:magnifier" size="16" />
                    </div>
                </div>
            </h6>

            <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-3 mb-1 text-muted">
                <span class="text-muted pb-2">博文分类</span>
            </h6>
            <ul class="nav flex-column">
                <li class="nav-item">
                    <a class="nav-link ps-4 bg-dark-subtle" aria-current="page" href="#">
                        Golang
                    </a>
                </li>
                <li class="nav-item">
                    <a class="nav-link ps-4" href="#">
                        TypeScript
                    </a>
                </li>
                <li class="nav-item">
                    <a class="nav-link ps-4" href="#">
                        MySQL
                    </a>
                </li>
            </ul>

            <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-3 mb-1 text-muted">
                <span class="text-muted pb-2">博文标签</span>
            </h6>
            <div class="nav mb-2 px-4">
                <button type="button" class="btn btn-primary btn-sm ms-2 mb-2" style="--bs-btn-padding-y:0.25rem;">Golang</button>
                <button type="button" class="btn btn-primary btn-sm ms-2 mb-2" style="--bs-btn-padding-y:0.25rem;">MySQL</button>
                <button type="button" class="btn btn-primary btn-sm ms-2 mb-2" style="--bs-btn-padding-y:0.25rem;">MongoDB</button>
                <button type="button" class="btn btn-primary btn-sm ms-2 mb-2" style="--bs-btn-padding-y:0.25rem;">My</button>
                <button type="button" class="btn btn-primary btn-sm ms-2 mb-2" style="--bs-btn-padding-y:0.25rem;">ABC</button>
                <button type="button" class="btn btn-primary btn-sm ms-2 mb-2" style="--bs-btn-padding-y:0.25rem;">MySQL</button>
            </div>
        </div>
    </aside>
</template>

<script setup lang="ts">

</script>

<style scoped lang="scss">
.my-input:focus {
    box-shadow: 0 0 0 0.1rem rgba(13, 110, 253, 0.25);
}
</style>