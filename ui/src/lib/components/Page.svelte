<script lang="ts">
	import { error_handling, pb } from "$lib/util";
	import { l, lang, iso } from "$lib/stores/lang";
	import type { Page } from "$lib/types/Page";
	import { find_page } from "$lib/stores/pages";

	export let path: string;

	let page: Page | null;
	$: {
		let right_path = path || location.pathname.substring(1);
		(async () => {
			page =
				find_page(path) ||
				(await pb
					.collection("page")
					.getFirstListItem(
						`path="${right_path}" && iso="${$iso}" && link!="deactivated"`
					)
					.catch((err) => {
						error_handling(err);
						return null;
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
