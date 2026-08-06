import type { Record } from "pocketbase";
import type { Product, ProductType } from "./Product";
import type { Order } from "./Order";
import type { Bar } from "./Bar";

export interface Serving extends Record {
	order: string; // relation
	product: string; // relation
	bar: string; // relation
	amount: number;
	price: number;
	free: boolean; // true for a wheel-win giveaway, not charged again
	type: ProductType; // client-side only, mirrors product.type for cart math, not a DB column
	expand: {
		order: Order;
		product: Product;
		bar: Bar;
	}
}
