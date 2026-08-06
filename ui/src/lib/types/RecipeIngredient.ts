import type { Record } from "pocketbase";
import type { Product } from "./Product";
import type { Ingredient } from "./Ingredient";

export interface RecipeIngredient extends Record {
	product: string; // relation
	ingredient: string; // relation
	quantity: number;
	expand: {
		product: Product;
		ingredient: Ingredient;
	}
}
