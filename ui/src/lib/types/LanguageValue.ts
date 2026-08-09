import type { Iso } from "$lib/types/Iso";
import type { RecordModel } from "pocketbase";

export interface LanguageValue extends RecordModel {
	language_key: string;
	iso: Iso;
	value: string;
}
