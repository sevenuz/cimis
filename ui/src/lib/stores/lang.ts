import { error_handling, pb } from "$lib/util";
import config from "$lib/config";
import type { LanguageValue } from "$lib/types/LanguageValue";
import type { LanguageKey } from "$lib/types/LanguageKey";
import { writable } from "svelte/store";
import type { Iso } from "$lib/types/Iso";

type LanguageDict = Record<string, Record<string, string>>;
export const lang = writable({} as LanguageDict);
export const iso = writable(config.default_language as Iso);

async function load_iso(iso: Iso) {
	lang.update(l => {
		l[iso] = {};
		return l;
	});
	const resultList: LanguageValue[] = await pb.collection('language_value').getFullList<LanguageValue>(100, {
		filter: `iso="${iso}"`,
		expand: 'language_key'
	}).catch(err => {
		error_handling(err);
		return [];
	});
	for (const lv of resultList) {
		lang.update(l => {
			l[iso][(lv.expand.language_key as LanguageKey).name] = lv.value;
			return l;
		});
	}
}

export function l(lang: LanguageDict, iso: Iso, key: string): string {
	if (Object.keys(lang).includes(iso)) {
		if (Object.keys(lang[iso]).includes(key)) {
			return lang[iso][key];
		}
	} else {
		load_iso(iso);
	}
	return key;
}
