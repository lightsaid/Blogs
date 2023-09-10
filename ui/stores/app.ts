import { defineStore } from "pinia"


type ThemeType = "dark" | "light" | "auto"

type AppType = {
    theme: ThemeType,
    isLogin: boolean
}

export const useAppStore = defineStore("app", {
    state: (): AppType => ({
        theme: "light",
        isLogin: false,
    }),

    actions: {
        setTheme(theme: ThemeType) {
            this.theme = theme
        },

        setLoginStatus(status: boolean){
            this.isLogin = status
        }
    },

})