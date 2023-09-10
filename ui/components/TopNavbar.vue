<template>
    <header class="topNavbar">
        <nav class="navbar navbar-expand-lg fixed-top navbar-dark bg-dark" style="--bs-navbar-padding-y:0.25rem;" aria-label="Main navigation">
            <div class="container">
                <NuxtLink class="navbar-brand fs-4" to="/">LightSaid's Blog</NuxtLink>
                <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav">
                    <span class="navbar-toggler-icon"></span>
                </button>

                <div class="collapse navbar-collapse" id="navbarNav">
                    <ul class="navbar-nav me-auto mb-2 mb-lg-0 fs-5">
                        <li class="nav-item ms-4">
                            <NuxtLink class="nav-link active" aria-current="page" href="/">Home</NuxtLink>
                        </li>
                        <li class="nav-item ms-4">
                            <NuxtLink class="nav-link" aria-current="page" href="/category">Category</NuxtLink>
                        </li>
                        <li class="nav-item ms-4">
                            <NuxtLink class="nav-link" aria-current="page" href="/tag">Tags</NuxtLink>
                        </li>
                        <li class="nav-item ms-4">
                            <NuxtLink class="nav-link" aria-current="page" href="/about">About</NuxtLink>
                        </li>
                    </ul>
                    <div class="d-flex">
                        <LoginDialog />
                        <RegisterDialog />
                    </div>

                    <div class="dropdown text-end ms-3">
                        <a href="#" class="d-block link-light text-decoration-none dropdown-toggle" id="dropdownUser1"
                            data-bs-toggle="dropdown" aria-expanded="false">
                            <img src="https://github.com/mdo.png" alt="mdo" width="32" height="32" class="rounded-circle">
                            LightSaid
                        </a>
                        <ul class="dropdown-menu text-small" aria-labelledby="dropdownUser1">
                            <li>
                                <NuxtLink class="dropdown-item" to="/about">设置</NuxtLink>
                            </li>
                            <li>
                                <NuxtLink class="dropdown-item" to="/posts/write">写博文</NuxtLink>
                            </li>
                            <li>
                                <hr class="dropdown-divider">
                            </li>
                            <li>
                                <NuxtLink class="dropdown-item" to="/">退出</NuxtLink>
                            </li>
                        </ul>
                    </div>

                    <div class="d-flex ms-3 align-items-center">
                        <a href="https://github.com/lightsaid" target="_blank">
                            <Icon name="uil:github" color="white" size="28" />
                        </a>
                        <button type="button" class="btn" 
                            data-bs-theme-value="light" 
                            v-show="appStore.theme === 'dark'" 
                            @click="appStore.setTheme('light')">
                            <Icon name="akar-icons:moon" color="white" size="24" data-bs-theme-value="light" @click="appStore.setTheme('light')"></Icon>
                        </button>
                        <button type="button" class="btn" 
                            data-bs-theme-value="dark" 
                            v-show="appStore.theme === 'light'"
                            @click="appStore.setTheme('dark')">
                            <Icon name="akar-icons:moon-fill" color="white" size="24" data-bs-theme-value="dark" @click="appStore.setTheme('dark')"></Icon>
                        </button>
                    </div>
                </div>
            </div>
        </nav>
    </header>
</template>

<script setup lang="ts">
import { Icon } from '#components'
import { useAppStore } from "~/stores/app"

const appStore = useAppStore()
console.log(appStore.theme)

</script>

<style scoped lang="scss">
.topNavbar {
    height: 56px;
    &>nav{
        // FIX： modal 显示时顶部栏抖动问题
        padding-right: 0 !important;
    }
}

.text-small {
  font-size: 90%;
}

.navbar,
.dropdown {
    transition: all 0.5s ease-in-out;
}

</style>