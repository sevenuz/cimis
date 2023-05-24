import { error_handling, pb, round_two_digits } from "$lib/util";
import type { Inventory } from "$lib/types/Inventory";
import { InventoryType } from "$lib/types/Inventory";
import type { PaymentMethod } from "$lib/types/PaymentMethod";
import { get, writable } from "svelte/store";
import type { Order } from "$lib/types/Order";
import type { Serving } from "$lib/types/Serving";
import { NotificationType, notify } from "./notifications";
import { iso, l, lang } from "./lang";

export const inventory = writable([] as Inventory[]);
export const payment_methods = writable([] as PaymentMethod[]);
export const order = writable({} as Order);
export const all_servings = writable([] as Serving[]);
export const new_servings = writable([] as Serving[]);

let loaded = false;

export async function load_bar() {
	if (loaded)
		return;
	loaded = true;
	const result_list = await pb
		.collection("inventory")
		.getFullList<Inventory>(50, {
			sort: "order",
			expand: "name"
		})
		.then((response) => {
			return response;
		})
		.catch((err) => {
			loaded = false;
			error_handling(err);
			return [] as Inventory[];
		});
	inventory.set(result_list);

	const result_list2 = await pb
		.collection("payment_method")
		.getFullList<PaymentMethod>(50, {
			expand: "name"
		})
		.then((response) => {
			return response;
		})
		.catch((err) => {
			loaded = false;
			error_handling(err);
			return [] as PaymentMethod[];
		});
	payment_methods.set(result_list2);

	const result_list3 = await pb
		.collection("serving")
		.getFullList<Serving>(50)
		.then((response) => {
			return response;
		})
		.catch((err) => {
			loaded = false;
			error_handling(err);
			return [] as Serving[];
		});
	all_servings.set(result_list3);
}

export function get_inventory_by_selection(selection: Serving): Inventory {
	return get(inventory).find(i => i.id == selection.product);
}

export function get_rest_inventory(product: Inventory): number {
	let rest = product.amount;
	for (let s of get(all_servings)) {
		if (s.product == product.id)
			rest -= s.amount;
	}
	for (let s of get(new_servings)) {
		if (s.product == product.id)
			rest -= s.amount;
	}
	return rest;
}

export function get_serving_total(selection: Serving[], pm: PaymentMethod): number {
	const percentages = [];
	let total = 0;
	for (let s of selection) {
		if (s.type == InventoryType.discount_percentage) {
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
	const product = get_inventory_by_selection(selection);
	const rest = get_rest_inventory(product);
	if (rest < 0) {
		notify({
			message: l(get(lang), get(iso), "ui_product_empty") + l(get(lang), get(iso), product.expand.name.name),
			type: NotificationType.error,
			duration: 2000
		});
		let p = get(new_servings).find(s => s.product == product.id)
		p.amount += rest;
	}
	if (selection.amount < 0 && product.type != InventoryType.deposit) {
		let p = get(new_servings).find(s => s.product == product.id)
		p.amount = 0; // this removes the entry
	}
	if (selection.amount > 1 && product.type == InventoryType.discount_percentage) {
		notify({
			message: l(get(lang), get(iso), "ui_discount_percentage_only_once"),
			type: NotificationType.warning,
			duration: 2000
		});
		let p = get(new_servings).find(s => s.product == product.id)
		p.amount = 1;
	}
	new_servings.set(get(new_servings).filter(s => s.amount != 0))
}

export function edit_selection(product: Inventory, n: number) {
	let p = get(new_servings).find(s => s.product == product.id)
	if (p) {
		p.amount += n;
		new_servings.set(get(new_servings));
	}
	else {
		p = {
			product: product.id,
			amount: n,
			price: product.price,
			type: product.type,
			expand: {
				product
			}
		} as Serving;
		new_servings.set([...get(new_servings), p])
	}
	check_selection(p);
}

export function remove_selection(selection: Serving) {
	new_servings.set(get(new_servings).filter(s => s != selection))
}
