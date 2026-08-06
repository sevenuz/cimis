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

	let remaining: { ingredient: Ingredient; amount: number }[] = [];
	let max_amount = 1;

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
		}));
		max_amount = Math.max(1, ...remaining.map((r) => Math.abs(r.amount)));
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
</script>

<div class="content text-center">
	<a class="rounded-full" href="/bar">&larr; {l($lang, $iso, "ui_back")}</a>
	<h1>{l($lang, $iso, "ui_inventory")}</h1>

	<div style="padding-bottom:20px;">
		{l($lang, $iso, "ui_event")}:
		<select class="rounded-full" bind:value={selected_event_id} on:change={compute}>
			{#each events as e}
				<option value={e.id}>{e.name}{e.active ? " *" : ""}</option>
			{/each}
		</select>
	</div>

	<div style="max-width:600px;margin:auto;">
		{#each remaining as r}
			<div class="text-left" style="margin-bottom:8px;">
				<div>{r.ingredient.name}: {r.amount} {r.ingredient.unit}</div>
				<div style="background:#3336; height:14px; border-radius:7px; overflow:hidden;">
					<div
						style="height:100%; width:{Math.max(
							0,
							(Math.abs(r.amount) / max_amount) * 100
						)}%; background:{r.amount < 0 ? '#e05252' : '#4caf50'};"
					/>
				</div>
			</div>
		{/each}
	</div>

	{#if $active_event && selected_event_id == $active_event.id}
		<div style="padding-top:30px;">
			<h2>{l($lang, $iso, "ui_add_inventory")}</h2>
			<select class="rounded-full" bind:value={new_ingredient_id}>
				<option value="">-</option>
				{#each $ingredients as i}
					<option value={i.id}>{i.name} ({i.unit})</option>
				{/each}
			</select>
			<input
				style="width: 100px;"
				class="rounded-full"
				type="number"
				bind:value={new_amount}
			/>
			<input class="rounded-full" bind:value={new_note} placeholder={l($lang, $iso, "ui_note")} />
			<button class="bg-white border-black rounded-full bg-yellow" on:click={add_receipt}>
				{l($lang, $iso, "ui_save")}
			</button>
		</div>
	{/if}
</div>

<svelte:head>
	<title>{l($lang, $iso, "ui_inventory")}</title>
</svelte:head>
