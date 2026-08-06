<script lang="ts">
	import { goto } from "$app/navigation";
	import { onMount } from "svelte";
	import { iso, l, lang } from "$lib/stores/lang";
	import { error_handling, pb } from "$lib/util";
	import type { User } from "$lib/types/User";
	import { Iso } from "$lib/types/Iso";
	import type { Product } from "$lib/types/Product";
	import { ProductType } from "$lib/types/Product";
	import type { RecipeIngredient } from "$lib/types/RecipeIngredient";
	import type { LanguageValue } from "$lib/types/LanguageValue";
	import { products, bars, ingredients, recipe_ingredients, load_catalog } from "$lib/stores/bar";
	import { NotificationType, notify } from "$lib/stores/notifications";

	let user: User | null = null;

	let editing: Product | null = null;
	let form_slug = "";
	let form_names: Record<string, string> = {};
	let form_price = 0;
	let form_type = ProductType.product;
	let form_color = "";
	let form_order = 0;
	let form_deactivated = false;
	let form_admin_only = false;
	let form_requires_deposit = false;
	let form_is_wheel = false;
	let form_instructions = "";
	let form_bars: string[] = [];

	let new_ingredient_id = "";
	let new_quantity = 0;

	let ingredient_name = "";
	let ingredient_unit = "";

	onMount(async () => {
		user = pb.authStore.model as User;
		if (!user || !user.admin) {
			goto("/bar");
			return;
		}
		await load_catalog(true);
	});

	function reset_form() {
		editing = null;
		form_slug = "";
		form_names = {};
		form_price = 0;
		form_type = ProductType.product;
		form_color = "";
		form_order = 0;
		form_deactivated = false;
		form_admin_only = false;
		form_requires_deposit = false;
		form_is_wheel = false;
		form_instructions = "";
		form_bars = [];
	}

	async function select_product(p: Product) {
		editing = p;
		form_slug = p.slug;
		form_price = p.price;
		form_type = p.type;
		form_color = p.color;
		form_order = p.order;
		form_deactivated = p.deactivated;
		form_admin_only = p.admin_only;
		form_requires_deposit = p.requires_deposit;
		form_is_wheel = p.is_wheel;
		form_instructions = p.instructions || "";
		form_bars = [...(p.bars || [])];

		form_names = {};
		const values = await pb
			.collection("language_value")
			.getFullList<LanguageValue>(undefined, { filter: `language_key="${p.name}"` })
			.catch(() => [] as LanguageValue[]);
		for (const v of values) form_names[v.iso] = v.value;
	}

	async function save_product() {
		if (!form_slug) return;

		let name_key_id = editing?.name;
		if (!name_key_id) {
			const key = await pb
				.collection("language_key")
				.create({ name: `bar_product_${form_slug}` })
				.catch(error_handling);
			if (!key) return;
			name_key_id = key.id;
		}

		for (const i of Object.values(Iso)) {
			const value = form_names[i];
			if (!value) continue;
			const existing = await pb
				.collection("language_value")
				.getFirstListItem<LanguageValue>(`language_key="${name_key_id}" && iso="${i}"`)
				.catch(() => null);
			if (existing) {
				await pb.collection("language_value").update(existing.id, { value }).catch(error_handling);
			} else {
				await pb
					.collection("language_value")
					.create({ language_key: name_key_id, iso: i, value })
					.catch(error_handling);
			}
		}

		const payload = {
			slug: form_slug,
			name: name_key_id,
			price: form_price,
			type: form_type,
			color: form_color,
			order: form_order,
			deactivated: form_deactivated,
			admin_only: form_admin_only,
			requires_deposit: form_requires_deposit,
			is_wheel: form_is_wheel,
			instructions: form_instructions,
			bars: form_bars,
		};

		if (editing) {
			await pb.collection("bar_product").update(editing.id, payload).catch(error_handling);
		} else {
			await pb.collection("bar_product").create(payload).catch(error_handling);
		}

		await load_catalog(true);
		notify({
			message: l($lang, $iso, "ui_order_saved"),
			type: NotificationType.success,
			duration: 2000,
		});
		reset_form();
	}

	function toggle_bar(bar_id: string) {
		form_bars = form_bars.includes(bar_id)
			? form_bars.filter((b) => b != bar_id)
			: [...form_bars, bar_id];
	}

	async function add_recipe_ingredient() {
		if (!editing || !new_ingredient_id) return;
		await pb
			.collection("bar_recipe_ingredient")
			.create({ product: editing.id, ingredient: new_ingredient_id, quantity: new_quantity })
			.catch(error_handling);
		new_ingredient_id = "";
		new_quantity = 0;
		await load_catalog(true);
	}

	async function remove_recipe_ingredient(ri: RecipeIngredient) {
		await pb.collection("bar_recipe_ingredient").delete(ri.id).catch(error_handling);
		await load_catalog(true);
	}

	async function add_ingredient() {
		if (!ingredient_name || !ingredient_unit) return;
		await pb
			.collection("bar_ingredient")
			.create({ name: ingredient_name, unit: ingredient_unit })
			.catch(error_handling);
		ingredient_name = "";
		ingredient_unit = "";
		await load_catalog(true);
	}

	$: editing_recipe_ingredients = editing
		? $recipe_ingredients.filter((ri) => ri.product == editing.id)
		: [];
</script>

<div class="content text-center">
	<a class="rounded-full" href="/bar">&larr; {l($lang, $iso, "ui_back")}</a>
	<h1>{l($lang, $iso, "ui_recipes")}</h1>

	<div class="md:grid grid-cols-2" style="width: 95vw; margin: auto; text-align: left;">
		<div>
			<h2>{l($lang, $iso, "ui_products")}</h2>
			<button class="rounded-full" on:click={reset_form}>+ {l($lang, $iso, "ui_new")}</button>
			<ul>
				{#each $products as p}
					<li>
						<button class="rounded-full" on:click={() => select_product(p)}>
							{l($lang, $iso, p.expand.name.name)} ({p.slug})
						</button>
					</li>
				{/each}
			</ul>
		</div>

		<div>
			<h2>{editing ? l($lang, $iso, "ui_edit") : l($lang, $iso, "ui_new")}</h2>
			<div>
				{l($lang, $iso, "ui_slug")}: <input class="rounded-full" bind:value={form_slug} />
			</div>
			{#each Object.values(Iso) as i}
				<div>
					{i}: <input class="rounded-full" bind:value={form_names[i]} />
				</div>
			{/each}
			<div>
				{l($lang, $iso, "ui_price")}:
				<input class="rounded-full" type="number" bind:value={form_price} />
			</div>
			<div>
				{l($lang, $iso, "ui_type")}:
				<select class="rounded-full" bind:value={form_type}>
					{#each Object.values(ProductType) as t}
						<option value={t}>{t}</option>
					{/each}
				</select>
			</div>
			<div>{l($lang, $iso, "ui_color")}: <input class="rounded-full" bind:value={form_color} /></div>
			<div>
				{l($lang, $iso, "ui_order")}:
				<input class="rounded-full" type="number" bind:value={form_order} />
			</div>
			<div>
				<label><input type="checkbox" bind:checked={form_deactivated} /> {l($lang, $iso, "ui_deactivated")}</label>
			</div>
			<div>
				<label><input type="checkbox" bind:checked={form_admin_only} /> {l($lang, $iso, "ui_admin_only")}</label>
			</div>
			<div>
				<label><input type="checkbox" bind:checked={form_requires_deposit} /> {l($lang, $iso, "ui_requires_deposit")}</label>
			</div>
			<div>
				<label><input type="checkbox" bind:checked={form_is_wheel} /> {l($lang, $iso, "ui_is_wheel")}</label>
			</div>
			<div>
				{l($lang, $iso, "ui_instructions")}:<br />
				<textarea class="rounded-md w-full" bind:value={form_instructions} />
			</div>
			<div>
				{l($lang, $iso, "ui_bar_selector")}:
				{#each $bars as b}
					<label>
						<input
							type="checkbox"
							checked={form_bars.includes(b.id)}
							on:change={() => toggle_bar(b.id)}
						/>
						{b.name}
					</label>
				{/each}
			</div>
			<button class="bg-white border-black rounded-full bg-yellow" on:click={save_product}>
				{l($lang, $iso, "ui_save")}
			</button>

			{#if editing}
				<h3 style="padding-top:20px;">{l($lang, $iso, "ui_recipe_ingredients")}</h3>
				<ul>
					{#each editing_recipe_ingredients as ri}
						<li>
							{ri.expand.ingredient.name}: {ri.quantity} {ri.expand.ingredient.unit}
							<button class="rounded-full" on:click={() => remove_recipe_ingredient(ri)}>x</button>
						</li>
					{/each}
				</ul>
				<select class="rounded-full" bind:value={new_ingredient_id}>
					<option value="">-</option>
					{#each $ingredients as i}
						<option value={i.id}>{i.name} ({i.unit})</option>
					{/each}
				</select>
				<input style="width: 100px;" class="rounded-full" type="number" bind:value={new_quantity} />
				<button class="rounded-full" on:click={add_recipe_ingredient}>+</button>
			{/if}
		</div>
	</div>

	<div style="padding-top:30px;">
		<h2>{l($lang, $iso, "ui_ingredients")}</h2>
		<ul>
			{#each $ingredients as i}
				<li>{i.name} ({i.unit})</li>
			{/each}
		</ul>
		<input class="rounded-full" bind:value={ingredient_name} placeholder={l($lang, $iso, "ui_name")} />
		<input class="rounded-full" bind:value={ingredient_unit} placeholder={l($lang, $iso, "ui_unit")} />
		<button class="rounded-full" on:click={add_ingredient}>+</button>
	</div>
</div>

<svelte:head>
	<title>{l($lang, $iso, "ui_recipes")}</title>
</svelte:head>
