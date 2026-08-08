<script lang="ts">
	import { createEventDispatcher } from "svelte";
	import { iso, l, lang, load_iso } from "$lib/stores/lang";
	import { error_handling, pb } from "$lib/util";
	import { Iso } from "$lib/types/Iso";
	import type { Product } from "$lib/types/Product";
	import { ProductType } from "$lib/types/Product";
	import type { RecipeIngredient } from "$lib/types/RecipeIngredient";
	import type { LanguageValue } from "$lib/types/LanguageValue";
	import { bars, ingredients, recipe_ingredients, load_catalog } from "$lib/stores/bar";
	import { NotificationType, notify } from "$lib/stores/notifications";

	// null = create mode. Once a create succeeds, this component keeps editing
	// the freshly made product in place (so recipe ingredients can be added
	// right away) without needing to navigate anywhere.
	export let product: Product | null = null;

	const dispatch = createEventDispatcher<{ saved: Product }>();

	let current: Product | null = null;
	let form_slug = "";
	let form_names: Record<string, string> = {};
	let form_price = 0;
	let form_type = ProductType.product;
	let form_color = "";
	let form_order = 0;
	let form_deactivated = false;
	let form_admin_only = false;
	let form_is_wheel = false;
	let form_instructions = "";
	let form_bars: string[] = [];

	let new_ingredient_id = "";
	let new_quantity = 0;

	// Tracks which `product` PROP id the form was last loaded from - not which
	// record was last saved. Never assign this outside this block: doing so
	// (e.g. after a create) makes it disagree with the still-unchanged `product`
	// prop, which retriggers this block and wipes `current` back to null via
	// load_form(product) right after a successful save.
	let last_loaded_id: string | null = null;
	$: if ((product?.id ?? null) !== last_loaded_id) {
		last_loaded_id = product?.id ?? null;
		load_form(product);
	}

	async function load_form(p: Product | null) {
		current = p;
		form_slug = p?.slug || "";
		form_price = p?.price || 0;
		form_type = p?.type || ProductType.product;
		form_color = p?.color || "";
		form_order = p?.order || 0;
		form_deactivated = p?.deactivated || false;
		form_admin_only = p?.admin_only || false;
		form_is_wheel = p?.is_wheel || false;
		form_instructions = p?.instructions || "";
		form_bars = [...(p?.bars || [])];

		form_names = {};
		if (!p) return;
		const values = await pb
			.collection("language_value")
			.getFullList<LanguageValue>(undefined, { filter: `language_key="${p.name}"` })
			.catch(() => [] as LanguageValue[]);
		for (const v of values) form_names[v.iso] = v.value;
	}

	async function save_product() {
		if (!form_slug) return;

		let name_key_id = current?.name;
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
			is_wheel: form_is_wheel,
			instructions: form_instructions,
			bars: form_bars,
		};

		let saved: Product | null;
		if (current) {
			saved = await pb
				.collection("bar_product")
				.update<Product>(current.id, payload)
				.catch((err) => {
					error_handling(err);
					return null;
				});
		} else {
			saved = await pb
				.collection("bar_product")
				.create<Product>(payload)
				.catch((err) => {
					error_handling(err);
					return null;
				});
		}
		if (!saved) return;
		current = saved;

		// l() caches translations per iso on first load and never notices new/changed
		// keys after that - force a refresh so the name we just saved shows up right away
		for (const i of Object.values(Iso)) {
			await load_iso(i);
		}

		await load_catalog(true);
		notify({
			message: l($lang, $iso, "ui_order_saved"),
			type: NotificationType.success,
			duration: 2000,
		});
		dispatch("saved", current);
	}

	function toggle_bar(bar_id: string) {
		form_bars = form_bars.includes(bar_id)
			? form_bars.filter((b) => b != bar_id)
			: [...form_bars, bar_id];
	}

	async function add_recipe_ingredient() {
		if (!current || !new_ingredient_id) return;
		await pb
			.collection("bar_recipe_ingredient")
			.create({ product: current.id, ingredient: new_ingredient_id, quantity: new_quantity })
			.catch(error_handling);
		new_ingredient_id = "";
		new_quantity = 0;
		await load_catalog(true);
	}

	async function update_recipe_ingredient(ri: RecipeIngredient, quantity: number) {
		await pb.collection("bar_recipe_ingredient").update(ri.id, { quantity }).catch(error_handling);
		await load_catalog(true);
	}

	async function remove_recipe_ingredient(ri: RecipeIngredient) {
		await pb.collection("bar_recipe_ingredient").delete(ri.id).catch(error_handling);
		await load_catalog(true);
	}

	$: editing_recipe_ingredients = current
		? $recipe_ingredients.filter((ri) => ri.product == current.id)
		: [];
</script>

<div style="max-width:700px; margin:auto; text-align:left;">
	<h2>Product</h2>

	<div class="form-grid">
		<label for="f-slug">{l($lang, $iso, "ui_slug")}</label>
		<input id="f-slug" class="form-input" bind:value={form_slug} />

		{#each Object.values(Iso) as i}
			<label for="f-name-{i}">{i}</label>
			<input id="f-name-{i}" class="form-input" bind:value={form_names[i]} />
		{/each}

		<label for="f-price">{l($lang, $iso, "ui_price")}</label>
		<input id="f-price" class="form-input" type="number" bind:value={form_price} />

		<label for="f-type">{l($lang, $iso, "ui_type")}</label>
		<select id="f-type" class="form-input" bind:value={form_type}>
			{#each Object.values(ProductType) as t}
				<option value={t}>{t}</option>
			{/each}
		</select>

		<label for="f-color">{l($lang, $iso, "ui_color")}</label>
		<input id="f-color" class="form-input" bind:value={form_color} />

		<label for="f-order">{l($lang, $iso, "ui_order")}</label>
		<input id="f-order" class="form-input" type="number" bind:value={form_order} />

		<label for="f-deactivated">{l($lang, $iso, "ui_deactivated")}</label>
		<input id="f-deactivated" type="checkbox" bind:checked={form_deactivated} />

		<label for="f-admin-only">{l($lang, $iso, "ui_admin_only")}</label>
		<input id="f-admin-only" type="checkbox" bind:checked={form_admin_only} />

		<label for="f-is-wheel">{l($lang, $iso, "ui_is_wheel")}</label>
		<input id="f-is-wheel" type="checkbox" bind:checked={form_is_wheel} />

		<label for="f-instructions">{l($lang, $iso, "ui_instructions")}</label>
		<textarea id="f-instructions" class="form-input" bind:value={form_instructions} />

		<span>{l($lang, $iso, "ui_bar_selector")}</span>
		<div>
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
	</div>

	<div class="text-center" style="padding-top:10px;">
		<button class="rounded-full" on:click={save_product}>
			{l($lang, $iso, "ui_save")}
		</button>
	</div>

	{#if current}
		<h3 style="padding-top:20px;">{l($lang, $iso, "ui_recipe_ingredients")}</h3>
		<ul>
			{#each editing_recipe_ingredients as ri}
				<li>
					{ri.expand.ingredient.name}:
					<input
						style="width: 80px;"
						class="form-input"
						type="number"
						value={ri.quantity}
						on:change={(e) => update_recipe_ingredient(ri, +e.currentTarget.value)}
					/>
					{ri.expand.ingredient.unit}
					<button class="rounded-full" on:click={() => remove_recipe_ingredient(ri)}>x</button>
				</li>
			{/each}
		</ul>
		<select class="form-input" bind:value={new_ingredient_id}>
			<option value="">-</option>
			{#each $ingredients as i}
				<option value={i.id}>{i.name} ({i.unit})</option>
			{/each}
		</select>
		<input style="width: 100px;" class="form-input" type="number" bind:value={new_quantity} />
		<button class="rounded-full" on:click={add_recipe_ingredient}>+</button>
	{/if}
</div>
