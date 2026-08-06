import { error_handling, local_cache_get, local_cache_set, pb, round_two_digits } from "$lib/util";
import type { Product } from "$lib/types/Product";
import { ProductType } from "$lib/types/Product";
import type { PaymentMethod } from "$lib/types/PaymentMethod";
import type { Bar } from "$lib/types/Bar";
import type { Event } from "$lib/types/Event";
import type { Ingredient } from "$lib/types/Ingredient";
import type { RecipeIngredient } from "$lib/types/RecipeIngredient";
import { get, writable } from "svelte/store";
import type { Order } from "$lib/types/Order";
import type { Serving } from "$lib/types/Serving";
import { NotificationType, notify } from "./notifications";
import { iso, l, lang } from "./lang";
import { enqueue, type QueuedOrder } from "./bar_queue";

const CURRENT_BAR_KEY = "bar_current_bar_id";

export const products = writable([] as Product[]);
export const payment_methods = writable([] as PaymentMethod[]);
export const bars = writable([] as Bar[]);
export const ingredients = writable([] as Ingredient[]);
export const recipe_ingredients = writable([] as RecipeIngredient[]);
export const active_event = writable(null as Event | null);
export const order = writable({} as Order);
export const new_servings = writable([] as Serving[]);
export const current_bar_id = writable(local_cache_get<string>(CURRENT_BAR_KEY) || "");
current_bar_id.subscribe((id) => local_cache_set(CURRENT_BAR_KEY, id));

async function load_with_cache<T>(
	collection: string,
	cache_key: string,
	params: Record<string, unknown>,
	set: (v: T[]) => void
) {
	try {
		const result = await pb.collection(collection).getFullList<T>(undefined, params);
		set(result);
		local_cache_set(cache_key, result);
	} catch (err) {
		const cached = local_cache_get<T[]>(cache_key);
		if (cached) {
			set(cached);
		} else {
			error_handling(err);
			set([]);
		}
	}
}

export async function load_orders(start_date: string, end_date: string, bar_id: string = "") {
	const filter = bar_id
		? `created >= "${start_date}" && created < "${end_date}" && bar="${bar_id}"`
		: `created >= "${start_date}" && created < "${end_date}"`;
	return await pb.collection("bar_order").getFullList<Order>(undefined, { filter });
}

export async function load_catalog(is_admin: boolean) {
	await load_with_cache<Product>(
		"bar_product",
		"cache_bar_product",
		{ filter: is_admin ? "" : "admin_only=false", sort: "order", expand: "name,bars" },
		(v) => products.set(v)
	);

	await load_with_cache<PaymentMethod>(
		"bar_payment_method",
		"cache_bar_payment_method",
		{ expand: "name" },
		(v) => payment_methods.set(v)
	);

	await load_with_cache<Bar>("bar", "cache_bar", {}, (v) => bars.set(v));

	await load_with_cache<Ingredient>(
		"bar_ingredient",
		"cache_bar_ingredient",
		{},
		(v) => ingredients.set(v)
	);

	await load_with_cache<RecipeIngredient>(
		"bar_recipe_ingredient",
		"cache_bar_recipe_ingredient",
		{ expand: "product,ingredient" },
		(v) => recipe_ingredients.set(v)
	);

	try {
		const ev = await pb.collection("event").getFirstListItem<Event>("active=true");
		active_event.set(ev);
		local_cache_set("cache_active_event", ev);
	} catch {
		active_event.set(local_cache_get<Event>("cache_active_event"));
	}
}

export function get_product_by_selection(selection: Serving): Product {
	return get(products).find((p) => p.id == selection.product);
}

export function get_serving_total(selection: Serving[], pm: PaymentMethod): number {
	const percentages = [];
	let total = 0;
	for (let s of selection) {
		if (s.type == ProductType.discount_percentage) {
			percentages.push(s.price);
		} else {
			total += s.amount * s.price;
		}
	}
	for (let p of percentages) {
		total *= p;
	}
	total *= pm.fee_percentage;
	total += pm.fee;
	return round_two_digits(total);
}

export function check_selection(selection: Serving) {
	const product = get_product_by_selection(selection);
	if (selection.amount < 0 && product.type != ProductType.deposit) {
		let p = get(new_servings).find((s) => s.product == product.id);
		p.amount = 0; // this removes the entry
	}
	if (selection.amount > 1 && product.type == ProductType.discount_percentage) {
		notify({
			message: l(get(lang), get(iso), "ui_discount_percentage_only_once"),
			type: NotificationType.warning,
			duration: 2000,
		});
		let p = get(new_servings).find((s) => s.product == product.id);
		p.amount = 1;
	}
	new_servings.set(get(new_servings).filter((s) => s.amount != 0));
}

export function edit_selection(product: Product, n: number) {
	let p = get(new_servings).find((s) => s.product == product.id);
	if (p) {
		p.amount += n;
		new_servings.set(get(new_servings));
	} else {
		p = {
			product: product.id,
			bar: get(current_bar_id),
			amount: n,
			price: product.price,
			free: false,
			type: product.type,
			expand: {
				product,
			},
		} as Serving;
		new_servings.set([...get(new_servings), p]);
	}
	check_selection(p);
}

export function remove_selection(selection: Serving) {
	new_servings.set(get(new_servings).filter((s) => s != selection));
}

// Adds a zero-priced serving line for a drink won on the lucky wheel - the wheel's
// fixed price was already charged as its own cart line, this just books the actual
// drink so its recipe ingredients get debited like a normal sale.
export function add_wheel_win(product: Product) {
	new_servings.set([
		...get(new_servings),
		{
			product: product.id,
			bar: get(current_bar_id),
			amount: 1,
			price: 0,
			free: true,
			type: product.type,
			expand: { product },
		} as Serving,
	]);
}

export function requires_deposit_warning(selection: Serving[]): boolean {
	return selection.some((s) => get_product_by_selection(s)?.requires_deposit);
}

export function queue_order(o: Order, servings: Serving[], bar_id?: string) {
	const bar = bar_id || get(current_bar_id);
	const item: QueuedOrder = {
		client_id: crypto.randomUUID(),
		user: o.user,
		payment_method: o.payment_method,
		bar,
		event: get(active_event)?.id || "",
		total: o.total,
		is_bookout: !!o.is_bookout,
		servings: servings.map((s) => ({
			product: s.product,
			bar: s.bar || bar,
			amount: s.amount,
			price: s.price,
			free: s.free,
		})),
	};
	enqueue(item);
}

export function queue_bookout(
	user_id: string,
	payment_method_id: string,
	amount: number,
	bar_id: string
) {
	queue_order(
		{
			user: user_id,
			payment_method: payment_method_id,
			total: amount,
			is_bookout: true,
		} as Order,
		[],
		bar_id
	);
}
