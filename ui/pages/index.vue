<template>
	<div class="container">
		<!-- 主栏 -->
		<div class="row">
			<div class="col-lg-9">
				<!-- 欢迎 -->
				<Welcome />

				<section class="posts mt-4">
					<template v-for="item in 8" :key="item">
						<PostCard />
					</template>
				</section>
			</div>

			<!-- 右侧栏 -->
			<div class="col-lg-3">
				<Sidebar />
			</div>
		</div>

	</div>
</template>

<script lang="ts" setup>

useHead({
  title: '首页',
  meta: [
    { name: 'description', content: '' }
  ],
  bodyAttrs: {
	style: 'background: url(pexels-pixabay-219906.jpg) no-repeat fixed center;',
  },
})



</script>

<style lang="scss" scoped></style>