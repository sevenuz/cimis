<script lang="ts">
	import { goto } from "$app/navigation";
	import { onMount } from "svelte";
	import { iso, l, lang } from "$lib/stores/lang";
	import { pb } from "$lib/util";
	import type { User } from "$lib/types/User";
	import type { Product } from "$lib/types/Product";
	import { load_catalog } from "$lib/stores/bar";
	import ProductEditor from "$lib/components/ProductEditor.svelte";

	let user: User | null = null;

	onMount(async () => {
		user = pb.authStore.model as User;
		if (!user || !user.admin) {
			goto("/bar");
			return;
		}
		await load_catalog(true);
	});

	function on_saved(e: CustomEvent<Product>) {
		goto(`/bar/products/edit/${e.detail.id}`);
	}
</script>

<div class="content text-center">
	<a class="nav-link-button" href="/bar/products">{l($lang, $iso, "ui_products")}</a>
	<h1>{l($lang, $iso, "ui_new")}</h1>

	<ProductEditor on:saved={on_saved} />
</div>

<svelte:head>
	<title>{l($lang, $iso, "ui_products")}</title>
</svelte:head>
