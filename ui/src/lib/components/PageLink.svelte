<script lang="ts">
	import { Addon, type Page } from "$lib/types/Page";

	export let page: Page;
</script>

{#if page.addon.includes(Addon.external_link)}
	<a href={page.path} class="text-white" rel="noreferrer" target="_blank">
		{page.title}
	</a>
{:else}
	<a href={"/" + page.path} class="text-white">
		{page.title}
	</a>
{/if}
