<script lang="ts">
	import { error_handling, pb } from "$lib/util";
	import { l, lang, iso } from "$lib/stores/lang";
	import type { Page } from "$lib/types/Page";
	import { find_page } from "$lib/stores/pages";

	export let path: string;

	async function load_page(p: string) {
		return pb
			.collection("page")
			.getFirstListItem(
				`path="${p}" && iso="${$iso}" && link!="deactivated"`
		)
	}

	let page: Page | null;
	$: {
		(async () => {
			page =
				find_page(path) ||
				(await load_page(path)
					.catch(async () => {
						return await load_page(location.pathname.substring(1))
							.catch((err) => {
								error_handling(err);
								return null;
							});
					}));
		})();
	}
</script>

<svelte:head>
	<title>{(page || {}).title || l($lang, $iso, "ui_default_page_title")}</title>
</svelte:head>

<div>
	{@html (page || {}).content || l($lang, $iso, "ui_404")}
</div>
