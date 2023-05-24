import PocketBase from 'pocketbase';
import config from "$lib/config";
import type { Record as PbRecord } from "pocketbase";
import { NotificationType, notify } from './stores/notifications';

export function isObject(obj: any) {
	return obj === Object(obj);
}

export const pb = new PocketBase(config.pb_url);
pb.autoCancellation(false); // TODO do we really want this?

const is_saving: Record<string, boolean> = {}; // checks if a request is already running for that survey_item
const update_ids: Record<string, string> = {}; // saves the id of a created record to call update next time
export async function upsert<T>(survey_item_id: string, data: T, collection_name: string = "survey_answer"): Promise<T & PbRecord> {
	const key = get_request_key(survey_item_id, collection_name);
	if (is_saving[key]) {
		return RejectFactory("Already saving " + key); // TODO: better error handling
	}

	is_saving[key] = true;
	let p: Promise<T & PbRecord>;
	let update_id = update_ids[key];

	if (update_id) {
		p = pb.collection(collection_name)
			.update<T & PbRecord>(update_id, data, { "$autoCancel": false })
	} else {
		p = pb.collection(collection_name)
			.create<T & PbRecord>(data, { "$autoCancel": false })
	}
	return p.then((response) => {
		add_update_id(survey_item_id, response.collectionName, response.id);
		return response;
	}).finally(() => {
		is_saving[key] = false;
	});
}

/*
* We calculate an id from item_id and collection_name because an item_id can
* sent multiple requests, for example for some nested items. Then the same item_id
* is used and collides when only this is used to identify the request...
* By using collection_name, a component can have one request for every collection
* without collision.
*/
function get_request_key(survey_item_id: string, collection_name: string): string {
	return survey_item_id + collection_name;
}

// use that function to set a loaded answer to use update in upsert fn
export function add_update_id(survey_item_id: string, collection_name: string, answer_id: string) {
	update_ids[get_request_key(survey_item_id, collection_name)] = answer_id;
}

// generic error handler for requests
export async function error_handling(error: Error) {
	console.log(error); // TODO: proper error handling
	notify({
		message: error.message,
		type: NotificationType.error,
		duration: 2000,
	});
}

export async function ResolveFactory<T>(data: T): Promise<T> {
	return new Promise((resolve, _) => resolve(data));
}

export async function RejectFactory<T>(error: string): Promise<T> {
	return new Promise((_, reject) => reject(new Error(error)));
}

export function round_two_digits(n: number): number {
	return Math.round(n * 100) / 100;
}
