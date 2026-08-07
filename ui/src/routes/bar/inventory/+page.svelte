<script lang="ts">
	import { goto } from "$app/navigation";
	import { onMount } from "svelte";
	import { iso, l, lang } from "$lib/stores/lang";
	import { error_handling, pb } from "$lib/util";
	import type { User } from "$lib/types/User";
	import type { Event } from "$lib/types/Event";
	import type { Ingredient } from "$lib/types/Ingredient";
	import type { Inventory } from "$lib/types/Inventory";
	import type { RecipeIngredient } from "$lib/types/RecipeIngredient";
	import {
		ingredients,
		recipe_ingredients,
		active_event,
		load_catalog,
	} from "$lib/stores/bar";
	import { NotificationType, notify } from "$lib/stores/notifications";

	let user: User | null = null;

	let events: Event[] = [];
	let selected_event_id = "";

	let remaining: { ingredient: Ingredient; amount: number; received: number }[] = [];

	let ingredient_name = "";
	let ingredient_unit = "";

	let new_ingredient_id = "";
	let new_amount = 0;
	let new_note = "";

	onMount(async () => {
		user = pb.authStore.model as User;
		if (!user || !user.admin) {
			goto("/bar");
			return;
		}
		await load_catalog(true);
		events = await pb.collection("event").getFullList<Event>().catch((err) => {
			error_handling(err);
			return [] as Event[];
		});
		selected_event_id = $active_event?.id || events[0]?.id || "";
		if (selected_event_id) compute();
	});

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

	async function add_receipt() {
		if (!new_ingredient_id || !$active_event) return;
		await pb
			.collection("bar_inventory")
			.create({
				event: $active_event.id,
				ingredient: new_ingredient_id,
				amount: new_amount,
				note: new_note,
			})
			.then(() => {
				new_amount = 0;
				new_note = "";
				notify({
					message: l($lang, $iso, "ui_order_saved"),
					type: NotificationType.success,
					duration: 2000,
				});
				if (selected_event_id == $active_event.id) compute();
			})
			.catch(error_handling);
	}

	async function compute() {
		if (!selected_event_id) return;

		const receipts = await pb
			.collection("bar_inventory")
			.getFullList<Inventory>(undefined, { filter: `event="${selected_event_id}"` })
			.catch((err) => {
				error_handling(err);
				return [] as Inventory[];
			});

		const servings = await pb
			.collection("bar_serving")
			.getFullList(undefined, { filter: `order.event="${selected_event_id}"` })
			.catch((err) => {
				error_handling(err);
				return [] as { product: string; amount: number }[];
			});

		const received: Record<string, number> = {};
		for (const r of receipts) {
			received[r.ingredient] = (received[r.ingredient] || 0) + r.amount;
		}

		const consumed: Record<string, number> = {};
		for (const s of servings) {
			for (const ri of $recipe_ingredients as RecipeIngredient[]) {
				if (ri.product == s.product) {
					consumed[ri.ingredient] = (consumed[ri.ingredient] || 0) + s.amount * ri.quantity;
				}
			}
		}

		remaining = $ingredients.map((i) => ({
			ingredient: i,
			amount: (received[i.id] || 0) - (consumed[i.id] || 0),
			received: received[i.id] || 0,
		}));
	}
</script>

<div class="content text-center">
	<a class="nav-link-button" href="/bar">&larr; {l($lang, $iso, "ui_back")}</a>
	<h1>{l($lang, $iso, "ui_inventory")}</h1>

	<div style="max-width:600px; margin:auto; text-align:left;">
		<details>
			<summary><h2 class="inline">{l($lang, $iso, "ui_ingredients")}</h2></summary>
			<div class="form-grid" style="padding-top:10px;">
				<label for="ing-name">{l($lang, $iso, "ui_name")}</label>
				<input id="ing-name" class="form-input" bind:value={ingredient_name} />
				<label for="ing-unit">{l($lang, $iso, "ui_unit")}</label>
				<input id="ing-unit" class="form-input" bind:value={ingredient_unit} />
			</div>
			<div class="text-center" style="padding-top:10px;">
				<button class="rounded-full" on:click={add_ingredient}>+</button>
			</div>
		</details>
	</div>

	{#if $active_event}
		<div style="max-width:600px; margin:auto; text-align:left; padding-top:10px;">
			<details>
				<summary><h2 class="inline">{l($lang, $iso, "ui_add_inventory")}</h2></summary>
				<div class="form-grid" style="padding-top:10px;">
					<label for="rec-ingredient">{l($lang, $iso, "ui_ingredients")}</label>
					<select id="rec-ingredient" class="form-input" bind:value={new_ingredient_id}>
						<option value="">-</option>
						{#each $ingredients as i}
							<option value={i.id}>{i.name} ({i.unit})</option>
						{/each}
					</select>
					<label for="rec-amount">{l($lang, $iso, "ui_amount")}</label>
					<input id="rec-amount" class="form-input" type="number" bind:value={new_amount} />
					<label for="rec-note">{l($lang, $iso, "ui_note")}</label>
					<input id="rec-note" class="form-input" bind:value={new_note} />
				</div>
				<div class="text-center" style="padding-top:10px;">
					<button class="bg-white border-black rounded-full bg-yellow" on:click={add_receipt}>
						{l($lang, $iso, "ui_save")}
					</button>
				</div>
			</details>
		</div>
	{/if}

	<div style="padding-top:30px;">
		{l($lang, $iso, "ui_event")}:
		<select class="form-input" bind:value={selected_event_id} on:change={compute}>
			{#each events as e}
				<option value={e.id}>{e.name}{e.active ? " *" : ""}</option>
			{/each}
		</select>
	</div>

	<div style="max-width:600px;margin:auto;padding-top:10px;">
		{#each remaining as r}
			<div class="text-left" style="margin-bottom:8px;">
				<div>{r.ingredient.name}: {r.amount} / {r.received} {r.ingredient.unit}</div>
				<div style="background:#3336; height:14px; border-radius:7px; overflow:hidden;">
					<div
						style="height:100%; width:{r.received > 0
							? Math.max(0, Math.min(100, (r.amount / r.received) * 100))
							: 0}%; background:{r.amount < 0 ? '#e05252' : '#4caf50'};"
					/>
				</div>
			</div>
		{/each}
	</div>
</div>

<svelte:head>
	<title>{l($lang, $iso, "ui_inventory")}</title>
</svelte:head>
