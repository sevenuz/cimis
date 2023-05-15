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
			{#each $navbarPages as page, id}
				{#if id !== 0}
					|
				{/if}
				<a href={"/" + page.path} class="hover:text-green">
					{l($lang, $iso, page.title)}
				</a>
			{/each}
		</div>
	</div>
</nav>
