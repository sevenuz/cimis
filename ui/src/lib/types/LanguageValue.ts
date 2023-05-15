import type { Iso } from "$lib/types/Iso";
import type { Record } from "pocketbase";

export interface LanguageValue extends Record {
	language_key: string;
	iso: Iso;
	value: string;
}
