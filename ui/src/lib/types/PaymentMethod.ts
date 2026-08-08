import type { Record } from "pocketbase";
import type { LanguageKey } from "./LanguageKey";

export interface PaymentMethod extends Record {
	name: string; // relation
	fee_percentage: number;
	fee: number;
	color: string;
	affects_box: boolean; // counts toward the physical cash box (e.g. cash, not card)
	expand: {
		name: LanguageKey
	}
}
