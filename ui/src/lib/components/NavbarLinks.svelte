<script lang="ts">
	import { l, lang, iso } from "$lib/stores/lang";
	import { onMount } from "svelte";
	import { navbarPages, load_pages } from "$lib/stores/pages";
    import PageLink from "./PageLink.svelte";

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
				<PageLink {page} />
			{/each}
		</div>
	</div>
</nav>
