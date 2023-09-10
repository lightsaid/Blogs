// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
	devtools: { enabled: true },
	modules: [
		'nuxt-icon',
		'@pinia/nuxt',
		'@pinia-plugin-persistedstate/nuxt',
	],
	css: [
		'@/assets/styles/index.scss'
	],
	devServer: {
		host: "0.0.0.0",
		port: 3000,
	},
})
