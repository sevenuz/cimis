<script lang="ts">
	import { goto } from "$app/navigation";
	import { onMount } from "svelte";
	import { iso, l, lang } from "$lib/stores/lang";
	import { pb, get_colors } from "$lib/util";
	import type { User } from "$lib/types/User";
	import { products, bars, recipe_ingredients, load_catalog } from "$lib/stores/bar";

	let user: User | null = null;

	onMount(async () => {
		user = pb.authStore.model as User;
		if (!user) {
			goto("/bar");
			return;
		}
		await load_catalog(user.admin);
	});

	$: recipe_ingredients_for = (product_id: string) =>
		$recipe_ingredients.filter((ri) => ri.product == product_id);

	$: bar_names_for = (bar_ids: string[]) =>
		bar_ids.map((id) => $bars.find((b) => b.id == id)?.name).filter(Boolean).join(", ");
</script>

<div class="content text-center">
	<a class="nav-link-button" href="/bar">{l($lang, $iso, "ui_back")}</a>
	<h1>{l($lang, $iso, "ui_products")}</h1>

	<div style="max-width:700px; margin:auto; text-align:left;">
		{#if user?.admin}
			<a class="nav-link-button" href="/bar/products/new">+ {l($lang, $iso, "ui_new")}</a>
		{/if}

		{#each $products as p}
			<details style="border-bottom: 1px solid rgb(222, 222, 222);">
				<summary>
					<b style={get_colors(p.color) + ";border-radius:2px;"}>{l($lang, $iso, p.expand.name.name)}</b> ({p.slug}) - <span class="text-white">{p.price}€</span> · {p.type}
					{#if p.bars?.length}· Bar: {bar_names_for(p.bars)}{/if}
					{#if p.deactivated}· {l($lang, $iso, "ui_deactivated")}{/if}
					{#if p.admin_only}· {l($lang, $iso, "ui_admin_only")}{/if}
					{#if p.is_wheel}· {l($lang, $iso, "ui_is_wheel")}{/if}
				</summary>
				<div style="padding: 10px 0;" class="text-white">
					{#if p.instructions}
						<p style="margin-bottom:8px;">{p.instructions}</p>
					{/if}
					<ul>
						{#each recipe_ingredients_for(p.id) as ri}
							<li>{ri.expand.ingredient.name}: {ri.quantity} {ri.expand.ingredient.unit}</li>
						{/each}
					</ul>
					{#if user?.admin}
						<a class="nav-link-button" href="/bar/products/edit/{p.id}">
							{l($lang, $iso, "ui_edit")}
						</a>
					{/if}
				</div>
			</details>
		{/each}
	</div>
</div>

<svelte:head>
	<title>{l($lang, $iso, "ui_products")}</title>
</svelte:head>
