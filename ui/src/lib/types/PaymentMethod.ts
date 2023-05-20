import type { Record } from "pocketbase";
import type { LanguageKey } from "./LanguageKey";

export interface PaymentMethod extends Record {
	name: string; // relation
	fee_percentage: number;
	fee: number;
	expand: {
		name: LanguageKey
	}
}
