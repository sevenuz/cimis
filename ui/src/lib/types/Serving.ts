import type { Record } from "pocketbase";
import type { Inventory } from "./Inventory";
import type { Order } from "./Order";

export interface Serving extends Record {
	order: string; // relation
	product: string; // relation
	amount: number;
	price: number;
	total: number;
	expand: {
		order: Order;
		product: Inventory;
	}
}
