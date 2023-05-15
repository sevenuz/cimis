<script lang="ts">
	import { l, lang, iso } from "$lib/stores/lang";
	import { onMount } from "svelte";
	import { navbarPages, load_pages } from "$lib/stores/pages";

	onMount(() => {
		load_pages();
	});
</script>

<nav>
	<div class="mt-5">
		<div>
			<a href={"/"} class="text-white">
				{l($lang, $iso, "ui_home_link")}
			</a>
			{#each $navbarPages as page}
				|
				<a href={"/" + page.path} class="text-white">
					{page.title}
				</a>
			{/each}
		</div>
	</div>
</nav>
