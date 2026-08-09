import type { RecordModel } from "pocketbase";
import type { Product } from "./Product";
import type { Ingredient } from "./Ingredient";

export interface RecipeIngredient extends RecordModel {
	product: string; // relation
	ingredient: string; // relation
	quantity: number;
	expand: {
		product: Product;
		ingredient: Ingredient;
	}
}
