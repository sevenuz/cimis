<script lang="ts">
	import { goto } from "$app/navigation";
	import { page } from "$app/stores";
	import { onMount } from "svelte";
	import { iso, l, lang } from "$lib/stores/lang";
	import { pb } from "$lib/util";
	import type { User } from "$lib/types/User";
	import { products, load_catalog } from "$lib/stores/bar";
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

	$: product = $products.find((p) => p.id === $page.params.id) ?? null;
</script>

<div class="content text-center">
	<a class="nav-link-button" href="/bar/products">{l($lang, $iso, "ui_products")}</a>
	<a class="nav-link-button" href="/bar/products/new">+ {l($lang, $iso, "ui_new")}</a>
	<h1>{l($lang, $iso, "ui_edit")}</h1>

	{#if product}
		<ProductEditor {product} />
	{/if}
</div>

<svelte:head>
	<title>{l($lang, $iso, "ui_products")}</title>
</svelte:head>
