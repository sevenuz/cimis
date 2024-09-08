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
export const new_servings = writable([] as Serving[]);

export async function load_orders(start_date: string, end_date: string) {
	return await pb.collection('bar_order').getFullList<Order>(50, {
		filter: `created >= "${start_date}" && created < "${end_date}"`,
	});
}

export async function load_bar(is_admin: boolean) {
	const result_list = await pb
		.collection("bar_inventory")
		.getFullList<Inventory>(50, {
			filter: is_admin ? "" : `admin_only=${is_admin}`,
			sort: "order",
			expand: "name"
		})
		.then((response) => {
			return response;
		})
		.catch((err) => {
			error_handling(err);
			return [] as Inventory[];
		});
	inventory.set(result_list);

	const result_list2 = await pb
		.collection("bar_payment_method")
		.getFullList<PaymentMethod>(50, {
			expand: "name"
		})
		.then((response) => {
			return response;
		})
		.catch((err) => {
			error_handling(err);
			return [] as PaymentMethod[];
		});
	payment_methods.set(result_list2);
}

export function get_inventory_by_selection(selection: Serving): Inventory {
	return get(inventory).find(i => i.id == selection.product);
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
