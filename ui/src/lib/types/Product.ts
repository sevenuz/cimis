import type { Record } from "pocketbase";
import type { LanguageKey } from "./LanguageKey";
import type { Bar } from "./Bar";

export enum ProductType {
	product = "product",
	discount_percentage = "discount_percentage",
	discount = "discount",
	deposit = "deposit"
}

export interface Product extends Record {
	slug: string;
	name: string; // relation
	price: number;
	type: ProductType;
	color: string;
	order: number;
	deactivated: boolean;
	admin_only: boolean;
	instructions: string;
	requires_deposit: boolean;
	is_wheel: boolean; // marks the one product that triggers the lucky-wheel win prompt
	bars: string[]; // relation, which bars sell this product
	expand: {
		name: LanguageKey;
		bars?: Bar[];
	}
}
